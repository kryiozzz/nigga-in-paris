package graph

import (
	"context"
	"database/sql"
	"fmt"
	"graphql/graph/model" // Ensure this path is correct
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// GetMyNotifications is the resolver for the getMyNotifications field.
func (r *queryResolver) GetMyNotifications(ctx context.Context, filter *string, limit *int32, offset *int32) ([]*model.Notification, error) {
	// 1. Get Current User ID
	currentUserID, err := getCurrentUserID(ctx) // From auth.go
	if err != nil {
		log.Printf("GetMyNotifications Error: Not authenticated: %v", err)
		// Return empty list if not authenticated, as they have no notifications
		return []*model.Notification{}, nil
		// Alternatively, return an error:
		// return nil, fmt.Errorf("authentication required")
	}

	// 2. Database Connection
	db, err := getDB() // From helper.go
	if err != nil {
		log.Printf("GetMyNotifications DB Error: %v", err)
		return nil, fmt.Errorf("internal server error connecting to DB")
	}
	defer db.Close()

	// 3. Build SQL Query
	var queryBuilder strings.Builder
	args := []interface{}{}
	argCounter := 1

	// Base query selecting necessary fields and joining accounts for triggering user info
	queryBuilder.WriteString(`
		SELECT
			n.notification_id, n.recipient_user_id, n.notification_type, n.entity_id, n.is_read, n.created_at,
			n.triggering_user_id,
			a.email, a.first_name, a.last_name, a.address, a.phone, a.age, a.gender, a.created_at as account_created_at, a.updated_at as account_updated_at
		FROM notifications n
		LEFT JOIN accounts a ON n.triggering_user_id = a.id
		WHERE n.recipient_user_id = $`)
	queryBuilder.WriteString(fmt.Sprintf("%d", argCounter))
	args = append(args, currentUserID)
	argCounter++

	// --- Apply Filtering ---
	if filter != nil {
		log.Printf("GetMyNotifications: Applying filter '%s'", *filter)
		if *filter == "unread" {
			// Filter only by unread status
			queryBuilder.WriteString(fmt.Sprintf(" AND n.is_read = false"))
		} else if *filter == "unread_new_post" {
			// *** ADDED: Filter by unread AND type 'new_post' ***
			queryBuilder.WriteString(fmt.Sprintf(" AND n.is_read = false AND n.notification_type = 'new_post'"))
		} else if *filter == "all" {
			// No additional filter needed if fetching all (read and unread)
		} else {
			// Optional: Handle unknown filters, maybe default to unread or log a warning
			log.Printf("GetMyNotifications: Unknown filter '%s' provided, defaulting to unread.", *filter)
			queryBuilder.WriteString(fmt.Sprintf(" AND n.is_read = false"))
		}
	} else {
		// Default filter if none is provided (e.g., only show unread)
		log.Println("GetMyNotifications: No filter provided, defaulting to unread.")
		queryBuilder.WriteString(fmt.Sprintf(" AND n.is_read = false"))
	}
	// --- End Filtering ---

	// Ordering
	queryBuilder.WriteString(" ORDER BY n.created_at DESC")

	// --- Apply Pagination ---
	actualLimit := int32(20) // Default limit
	if limit != nil && *limit > 0 {
		// Apply a reasonable maximum limit if desired
		// if *limit > 100 { actualLimit = 100 } else { actualLimit = *limit }
		actualLimit = *limit
	}
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argCounter))
	args = append(args, actualLimit)
	argCounter++

	actualOffset := int32(0) // Default offset
	if offset != nil && *offset >= 0 {
		actualOffset = *offset
	}
	queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", argCounter))
	args = append(args, actualOffset)
	argCounter++
	// --- End Pagination ---

	finalQuery := queryBuilder.String()
	log.Printf("Executing GetMyNotifications query: [%s] with args: %v", finalQuery, args)

	// 4. Execute Query
	queryCtx, cancelQuery := context.WithTimeout(ctx, 10*time.Second)
	defer cancelQuery()
	rows, err := db.QueryContext(queryCtx, finalQuery, args...)
	if err != nil {
		log.Printf("GetMyNotifications DB Error executing query: %v", err)
		return nil, fmt.Errorf("failed to fetch notifications")
	}
	defer rows.Close()

	// 5. Scan Results
	notifications := []*model.Notification{}
	for rows.Next() {
		var notif model.Notification
		var triggeringUserID sql.NullString
		var entityID sql.NullString
		var createdAt time.Time
		var accEmail, accFirstName, accLastName, accAddress, accPhone, accGender sql.NullString
		var accAge sql.NullInt32
		var accCreatedAt, accUpdatedAt sql.NullTime

		err := rows.Scan(
			&notif.NotificationID, &notif.RecipientUserID, &notif.NotificationType, &entityID, &notif.IsRead, &createdAt,
			&triggeringUserID,
			&accEmail, &accFirstName, &accLastName, &accAddress, &accPhone, &accAge, &accGender, &accCreatedAt, &accUpdatedAt,
		)
		if err != nil {
			log.Printf("GetMyNotifications DB Error scanning row: %v", err)
			continue // Skip this notification
		}

		notif.CreatedAt = createdAt.Format(time.RFC3339)
		if entityID.Valid {
			notif.EntityID = &entityID.String
		} else {
			notif.EntityID = nil
		}

		if triggeringUserID.Valid {
			triggeringUser := &model.Account{AccountID: triggeringUserID.String}
			if accEmail.Valid {
				triggeringUser.Email = accEmail.String
			}
			if accFirstName.Valid {
				triggeringUser.FirstName = accFirstName.String
			}
			if accLastName.Valid {
				triggeringUser.LastName = accLastName.String
			}
			if accAddress.Valid {
				triggeringUser.Address = &accAddress.String
			}
			if accPhone.Valid {
				triggeringUser.Phone = &accPhone.String
			}
			if accAge.Valid {
				triggeringUser.Age = accAge.Int32
			}
			if accGender.Valid {
				triggeringUser.Gender = &accGender.String
			}
			if accCreatedAt.Valid {
				triggeringUser.CreatedAt = accCreatedAt.Time.Format(time.RFC3339)
			}
			if accUpdatedAt.Valid {
				updatedAtStr := accUpdatedAt.Time.Format(time.RFC3339)
				triggeringUser.UpdatedAt = &updatedAtStr
			}
			notif.TriggeringUser = triggeringUser
		} else {
			notif.TriggeringUser = nil
		}
		notifications = append(notifications, &notif)
	}

	// 6. Check for errors during row iteration
	if err = rows.Err(); err != nil {
		log.Printf("GetMyNotifications DB Error iterating rows: %v", err)
		return nil, fmt.Errorf("error reading notifications list")
	}

	// 7. Return
	log.Printf("GetMyNotifications: Returning %d notifications for user %s with filter '%v'", len(notifications), currentUserID, filter)
	return notifications, nil
} // End of GetMyNotifications

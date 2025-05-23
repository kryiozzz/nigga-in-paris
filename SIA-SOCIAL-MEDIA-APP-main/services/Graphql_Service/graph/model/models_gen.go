// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Account struct {
	AccountID   string  `json:"accountId"`
	Email       string  `json:"email"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Address     *string `json:"address,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Age         int32   `json:"age"`
	Gender      *string `json:"gender,omitempty"`
	IsFollowing *bool   `json:"isFollowing,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
}

type CreatePostInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID string `json:"authorId"`
}

type CreateProfileInput struct {
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	FirstName         *string `json:"firstName,omitempty"`
	MiddleName        *string `json:"middleName,omitempty"`
	LastName          *string `json:"lastName,omitempty"`
	Bio               *string `json:"bio,omitempty"`
	ProfilePictureURL *string `json:"profilePictureUrl,omitempty"`
	BannerPictureURL  *string `json:"bannerPictureUrl,omitempty"`
	DateOfBirth       *string `json:"dateOfBirth,omitempty"`
	Address           *string `json:"address,omitempty"`
}

type Mutation struct {
}

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

// Represents a notification for a user.
type Notification struct {
	NotificationID   string   `json:"notificationId"`
	RecipientUserID  string   `json:"recipientUserId"`
	TriggeringUser   *Account `json:"triggeringUser,omitempty"`
	NotificationType string   `json:"notificationType"`
	EntityID         *string  `json:"entityId,omitempty"`
	IsRead           bool     `json:"isRead"`
	CreatedAt        string   `json:"createdAt"`
	// The post associated with this notification, if applicable (e.g., for 'new_post', 'like', 'new_comment').
	Post *Post `json:"post,omitempty"`
}

type Post struct {
	PostID    string   `json:"postId"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	AuthorID  string   `json:"authorId"`
	Author    *Account `json:"author"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt *string  `json:"updatedAt,omitempty"`
}

type Profile struct {
	ProfileID         string  `json:"profileId"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Password          string  `json:"password"`
	FirstName         *string `json:"firstName,omitempty"`
	MiddleName        *string `json:"middleName,omitempty"`
	LastName          *string `json:"lastName,omitempty"`
	Bio               *string `json:"bio,omitempty"`
	ProfilePictureURL *string `json:"profilePictureUrl,omitempty"`
	BannerPictureURL  *string `json:"bannerPictureUrl,omitempty"`
	DateOfBirth       *string `json:"dateOfBirth,omitempty"`
	Address           *string `json:"address,omitempty"`
}

type Query struct {
}

type RegisterInput struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Address   *string `json:"address,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Age       int32   `json:"age"`
	Gender    *string `json:"gender,omitempty"`
}

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

package linkfy

import (
	"time"

	"github.com/google/uuid"
)

// Linkfy represents a linkfy profile in the database
type Linkfy struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
	Name      string    `json:"name"`
	Bio       *string   `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LinkfyCreated struct {
	UserID    string  `json:"user_id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Name      string  `json:"name"`
	Bio       *string `json:"bio"`
}

type LinkfyUpdated struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
	Name      string    `json:"name"`
	Bio       *string   `json:"bio"`
}

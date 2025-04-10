package linkfylink

import (
	"time"

	"github.com/google/uuid"
)

// LinkfyLink represents a link in the linkfy profile
type LinkfyLink struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	NameURL   string    `json:"name_url"`
	IconsURL  string    `json:"icons_url"`
	CreatedAt time.Time `json:"created_at"`
}

// LinkfyLinkCreated represents data needed for creating a new linkfy link
type LinkfyLinkCreated struct {
	Name     string `json:"name"`
	NameURL  string `json:"name_url"`
	IconsURL string `json:"icons_url"`
}

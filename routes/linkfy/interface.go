package linkfy

import (
	"github.com/google/uuid"
	"github.com/yantology/linkfy/pkg/customerror"
)

type LinkfyDBInterface interface {
	// CreateLinkfy creates a new linkfy profile
	CreateLinkfy(linkfy *LinkfyCreated) *customerror.CustomError

	// GetLinkfyByID retrieves a linkfy profile by its ID
	GetLinkfyByID(id uuid.UUID) (*Linkfy, *customerror.CustomError)

	// GetLinkfyByUsername retrieves a linkfy profile by username
	GetLinkfyByUsername(username string) (*Linkfy, *customerror.CustomError)

	// GetAllLinkfyByUserID retrieves all linkfy profiles for a specific user
	GetAllLinkfyByUserID(userID string) ([]*Linkfy, *customerror.CustomError)

	// UpdateLinkfy updates an existing linkfy profile
	UpdateLinkfy(linkfy *LinkfyUpdated) *customerror.CustomError

	// CheckUsernameExists checks if a username already exists (for debounce)
	CheckUsernameExists(username string) *customerror.CustomError

	CheckUsernameNotExists(username string) *customerror.CustomError

	// DeleteLinkfy deletes a linkfy profile by its ID
	DeleteLinkfy(id uuid.UUID, userID string) *customerror.CustomError
}

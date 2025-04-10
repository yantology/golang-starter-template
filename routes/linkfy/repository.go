package linkfy

import (
	"github.com/google/uuid"
	"github.com/yantology/linkfy/pkg/customerror"
)

type linkfyRepository struct {
	db LinkfyDBInterface
}

// NewLinkfyRepository creates a new instance of the linkfy repository
func NewLinkfyRepository(db LinkfyDBInterface) *linkfyRepository {
	return &linkfyRepository{db: db}
}

// CreateLinkfy creates a new linkfy profile
func (r *linkfyRepository) CreateLinkfy(linkfyCreated *LinkfyCreated) *customerror.CustomError {
	return r.db.CreateLinkfy(linkfyCreated)
}

// GetLinkfyByID retrieves a linkfy profile by its ID
func (r *linkfyRepository) GetLinkfyByID(id uuid.UUID) (*Linkfy, *customerror.CustomError) {
	return r.db.GetLinkfyByID(id)
}

// GetLinkfyByUsername retrieves a linkfy profile by username
func (r *linkfyRepository) GetLinkfyByUsername(username string) (*Linkfy, *customerror.CustomError) {
	return r.db.GetLinkfyByUsername(username)
}

// GetAllLinkfyByUserID retrieves all linkfy profiles for a specific user
func (r *linkfyRepository) GetAllLinkfyByUserID(userID string) ([]*Linkfy, *customerror.CustomError) {
	return r.db.GetAllLinkfyByUserID(userID)
}

// UpdateLinkfy updates an existing linkfy profile
func (r *linkfyRepository) UpdateLinkfy(linkfyUpdate *LinkfyUpdated) *customerror.CustomError {
	return r.db.UpdateLinkfy(linkfyUpdate)
}

// CheckUsernameExists checks if a username already exists (for debounce)
func (r *linkfyRepository) CheckUsernameExists(username string) *customerror.CustomError {
	return r.db.CheckUsernameExists(username)
}

func (r *linkfyRepository) CheckUsernameNotExists(username string) *customerror.CustomError {
	return r.db.CheckUsernameNotExists(username)
}

// DeleteLinkfy deletes a linkfy profile by its ID
func (r *linkfyRepository) DeleteLinkfy(id uuid.UUID, userID string) *customerror.CustomError {
	return r.db.DeleteLinkfy(id, userID)
}

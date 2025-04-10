package linkfylink

import (
	"github.com/yantology/linkfy/pkg/customerror"
)

type linkfyLinkRepository struct {
	db LinkfyLinkDBInterface
}

// NewLinkfyLinkRepository creates a new instance of the linkfy link repository
func NewLinkfyLinkRepository(db LinkfyLinkDBInterface) *linkfyLinkRepository {
	return &linkfyLinkRepository{db: db}
}

// CreateLinks creates multiple links for a linkfy profile with specified order
func (r *linkfyLinkRepository) CreateLinks(linkfy_id string, links []*LinkfyLinkCreated) *customerror.CustomError {
	return r.db.CreateLinks(linkfy_id, links)
}

// GetLinkByLinkfyID retrieves a link by its ID
func (r *linkfyLinkRepository) GetLinkByLinkfyID(linkfy_id string) ([]*LinkfyLink, *customerror.CustomError) {
	return r.db.GetLinkByLinkfyID(linkfy_id)
}

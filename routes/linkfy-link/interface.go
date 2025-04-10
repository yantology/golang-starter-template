package linkfylink

import (
	"github.com/yantology/linkfy/pkg/customerror"
)

type LinkfyLinkDBInterface interface {
	// CreateLinks creates multiple links for a linkfy profile with specified order
	CreateLinks(linkfy_id string, links []*LinkfyLinkCreated) *customerror.CustomError

	// GetLinkByLinkfyID retrieves a link by its ID
	GetLinkByLinkfyID(linkfy_id string) ([]*LinkfyLink, *customerror.CustomError)
}

package linkfylink

// Request DTOs

// CreateLinksRequest represents the request for creating new links for a linkfy profile
// @Description Create links request model
type CreateLinksRequest struct {
	Links []LinkRequest `json:"links" binding:"required"`
}

// LinkRequest represents a single link in the create links request
// @Description Link request model
type LinkRequest struct {
	Name     string `json:"name" binding:"required" example:"GitHub"`
	NameURL  string `json:"name_url" binding:"required" example:"https://github.com/username"`
	IconsURL string `json:"icons_url" binding:"required" example:"https://example.com/github-icon.png"`
}

// Response DTOs

// MessageResponse represents a generic message response
// @Description Generic message response model
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// DataResponse represents a generic data response
// @Description Generic data response model
type DataResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message" example:"Operation completed successfully"`
}

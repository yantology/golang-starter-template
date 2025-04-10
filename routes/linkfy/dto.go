package linkfy

// Request DTOs

// CreateLinkfyRequest represents the request for creating a new linkfy profile
// @Description Create Linkfy profile request model
type CreateLinkfyRequest struct {
	Username  string  `json:"username" binding:"required" example:"johndoe"`
	Name      string  `json:"name" binding:"required" example:"John Doe"`
	AvatarURL *string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Bio       *string `json:"bio" example:"Web developer passionate about Go"`
}

// UpdateLinkfyRequest represents the request for updating a linkfy profile
// @Description Update Linkfy profile request model
type UpdateLinkfyRequest struct {
	Username  string  `json:"username" example:"johndoe"`
	Name      string  `json:"name" example:"John Doe"`
	AvatarURL *string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Bio       *string `json:"bio" example:"Web developer passionate about Go"`
}

// CheckUsernameRequest represents the request for checking username availability
// @Description Check username availability request model
type CheckUsernameRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
}

// UsernameExistsResponse represents the response for the username check
// @Description Username availability response model
type UsernameExistsResponse struct {
	Exists bool `json:"exists" example:"false"`
}

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

package linkfy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yantology/linkfy/middleware"
)

type LinkfyHandler struct {
	service    *LinkfyService
	repository *linkfyRepository
}

func NewLinkfyHandler(service *LinkfyService, repository *linkfyRepository) *LinkfyHandler {
	return &LinkfyHandler{
		service:    service,
		repository: repository,
	}
}

// RegisterRoutes registers all the routes for linkfy
func (h *LinkfyHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", h.CreateLinkfy)
	router.GET("/:linkfy_id", h.GetLinkfyByID)
	router.GET("/username/:username", h.GetLinkfyByUsername)
	router.GET("", h.GetAllLinkfy)
	router.PUT("/:linkfy_id", h.UpdateLinkfy)
	router.POST("/check-username", h.CheckUsername)
	router.DELETE("/:linkfy_id", h.DeleteLinkfy)
}

// CreateLinkfy handles the creation of a new linkfy profile
// @Summary Create a new linkfy profile
// @Description Create a new linkfy profile for the authenticated user
// @Tags linkfy
// @Accept json
// @Produce json
// @Param request body CreateLinkfyRequest true "Linkfy profile details"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} MessageResponse
// @Failure 401 {object} MessageResponse
// @Failure 409 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy [post]
func (h *LinkfyHandler) CreateLinkfy(c *gin.Context) {
	// Get user ID from context
	userClaims := middleware.ExtractUserClaims(c)
	if userClaims == nil || userClaims.UserID == "" {
		c.JSON(http.StatusUnauthorized, MessageResponse{
			Message: "User tidak terautentikasi",
		})
		return
	}

	var req CreateLinkfyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Sanitize username
	if custErr := h.service.UsernameSanitize(req.Username); custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	// Create linkfy profile
	avatar := req.AvatarURL
	if avatar == nil {
		defaultAvatar := "https://github.com/shadcn.png"
		avatar = &defaultAvatar
	}
	bio := req.Bio
	if bio == nil {
		defaultBio := "No bio provided"
		bio = &defaultBio
	}
	LinkfyCreated := &LinkfyCreated{
		UserID:    userClaims.UserID,
		Username:  req.Username,
		Name:      req.Name,
		AvatarURL: avatar,
		Bio:       bio,
	}

	if custErr := h.repository.CreateLinkfy(LinkfyCreated); custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{
		Message: "Linkfy profile created successfully",
	})

}

// GetLinkfyByID handles getting a linkfy profile by ID
// @Summary Get linkfy profile by ID
// @Description Get a specific linkfy profile by its ID
// @Tags linkfy
// @Accept json
// @Produce json
// @Param linkfy_id path string true "Linkfy profile ID"
// @Success 200 {object} DataResponse[Linkfy]
// @Failure 400 {object} MessageResponse
// @Failure 401 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/{linkfy_id} [get]
func (h *LinkfyHandler) GetLinkfyByID(c *gin.Context) {
	idStr := c.Param("linkfy_id")
	linkfy_id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid ID format",
		})
		return
	}

	linkfy, custErr := h.repository.GetLinkfyByID(linkfy_id)
	if custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse[*Linkfy]{
		Data:    linkfy,
		Message: "Linkfy profile retrieved successfully",
	})
}

// GetLinkfyByUsername handles getting a linkfy profile by username
// @Summary Get linkfy profile by username
// @Description Get a specific linkfy profile by its username
// @Tags linkfy
// @Accept json
// @Produce json
// @Param username path string true "Linkfy username"
// @Success 200 {object} DataResponse[Linkfy]
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/username/{username} [get]
func (h *LinkfyHandler) GetLinkfyByUsername(c *gin.Context) {
	username := c.Param("username")

	linkfy, custErr := h.repository.GetLinkfyByUsername(username)
	if custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse[*Linkfy]{
		Data:    linkfy,
		Message: "Linkfy profile retrieved successfully",
	})
}

// GetAllLinkfy handles getting all linkfy profiles for the authenticated user
// @Summary Get all user's linkfy profiles
// @Description Get all linkfy profiles for the authenticated user
// @Tags linkfy
// @Accept json
// @Produce json
// @Success 200 {object} DataResponse[[]Linkfy]
// @Failure 401 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy [get]
func (h *LinkfyHandler) GetAllLinkfy(c *gin.Context) {
	userClaims := middleware.ExtractUserClaims(c)
	if userClaims == nil || userClaims.UserID == "" {
		c.JSON(http.StatusUnauthorized, MessageResponse{
			Message: "User tidak terautentikasi",
		})
		return
	}

	linkfies, custErr := h.repository.GetAllLinkfyByUserID(userClaims.UserID)
	if custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse[[](*Linkfy)]{
		Data:    linkfies,
		Message: "Linkfy profiles retrieved successfully",
	})
}

// UpdateLinkfy handles updating a linkfy profile
// @Summary Update a linkfy profile
// @Description Update an existing linkfy profile for the authenticated user
// @Tags linkfy
// @Accept json
// @Produce json
// @Param linkfy_id path string true "Linkfy profile ID"
// @Param request body UpdateLinkfyRequest true "Updated linkfy profile details"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} MessageResponse
// @Failure 401 {object} MessageResponse
// @Failure 403 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Failure 409 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/{linkfy_id} [put]
func (h *LinkfyHandler) UpdateLinkfy(c *gin.Context) {
	userClaims := middleware.ExtractUserClaims(c)
	if userClaims == nil || userClaims.UserID == "" {
		c.JSON(http.StatusUnauthorized, MessageResponse{
			Message: "User tidak terautentikasi",
		})
		return
	}

	idStr := c.Param("linkfy_id")
	linkfy_id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid ID format",
		})
		return
	}

	var req UpdateLinkfyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Sanitize username
	if custErr := h.service.UsernameSanitize(req.Username); custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}
	// Create linkfy profile
	avatar := req.AvatarURL
	if avatar == nil {
		defaultAvatar := "https://github.com/shadcn.png"
		avatar = &defaultAvatar
	}
	bio := req.Bio
	if bio == nil {
		defaultBio := "No bio provided"
		bio = &defaultBio
	}
	LinkfyUpdated := &LinkfyUpdated{
		ID:        linkfy_id,
		UserID:    userClaims.UserID,
		Username:  req.Username,
		Name:      req.Name,
		AvatarURL: avatar,
		Bio:       bio,
	}

	if custErr := h.repository.UpdateLinkfy(LinkfyUpdated); custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Linkfy profile updated successfully",
	})
}

// CheckUsername handles checking if a username is available
// @Summary Check username availability
// @Description Check if a username is available for registration
// @Tags linkfy
// @Accept json
// @Produce json
// @Param request body CheckUsernameRequest true "Username to check"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/check-username [post]
func (h *LinkfyHandler) CheckUsername(c *gin.Context) {
	var req CheckUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	custErr := h.repository.CheckUsernameNotExists(req.Username)
	if custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Username is available",
	})
}

// DeleteLinkfy handles deleting a linkfy profile
// @Summary Delete a linkfy profile
// @Description Delete an existing linkfy profile for the authenticated user
// @Tags linkfy
// @Accept json
// @Produce json
// @Param linkfy_id path string true "Linkfy profile ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} MessageResponse
// @Failure 401 {object} MessageResponse
// @Failure 403 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/{linkfy_id} [delete]
func (h *LinkfyHandler) DeleteLinkfy(c *gin.Context) {
	userClaims := middleware.ExtractUserClaims(c)
	if userClaims == nil || userClaims.UserID == "" {
		c.JSON(http.StatusUnauthorized, MessageResponse{
			Message: "User tidak terautentikasi",
		})
		return
	}

	idStr := c.Param("linkfy_id")
	linkfy_id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid ID format",
		})
		return
	}

	custErr := h.repository.DeleteLinkfy(linkfy_id, userClaims.UserID)

	if custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Linkfy profile deleted successfully",
	})
}

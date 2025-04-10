package linkfylink

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yantology/linkfy/middleware"
)

type LinkfyLinkHandler struct {
	repository *linkfyLinkRepository
}

func NewLinkfyLinkHandler(repository *linkfyLinkRepository) *LinkfyLinkHandler {
	return &LinkfyLinkHandler{
		repository: repository,
	}
}

// RegisterRoutes registers all the routes for linkfy-link
func (h *LinkfyLinkHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/:linkfy_id/links", h.CreateLinks)
	router.GET("/:linkfy_id/links", h.GetLinksByLinkfyID)
}

// CreateLinks handles the creation of new links for a linkfy profile
// @Summary Create new links for a linkfy profile
// @Description Create multiple links for a specified linkfy profile
// @Tags linkfy-link
// @Accept json
// @Produce json
// @Param linkfy_id path string true "Linkfy profile ID"
// @Param request body CreateLinksRequest true "Links details"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} MessageResponse
// @Failure 401 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/{linkfy_id}/links [post]
func (h *LinkfyLinkHandler) CreateLinks(c *gin.Context) {
	// Get user ID from context
	userClaims := middleware.ExtractUserClaims(c)
	if userClaims == nil || userClaims.UserID == "" {
		c.JSON(http.StatusUnauthorized, MessageResponse{
			Message: "User tidak terautentikasi",
		})
		return
	}

	// Get linkfy_id from path
	linkfyIDStr := c.Param("linkfy_id")
	if linkfyIDStr == "" {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Linkfy ID is required",
		})
		return
	}

	var req CreateLinksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert request to model
	links := make([]*LinkfyLinkCreated, 0, len(req.Links))
	for _, link := range req.Links {
		links = append(links, &LinkfyLinkCreated{
			Name:     link.Name,
			NameURL:  link.NameURL,
			IconsURL: link.IconsURL,
		})
	}

	// Create the links
	if custErr := h.repository.CreateLinks(linkfyIDStr, links); custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{
		Message: "Links created successfully",
	})
}

// GetLinksByLinkfyID handles retrieving all links for a linkfy profile
// @Summary Get links by linkfy ID
// @Description Get all links associated with a specific linkfy profile
// @Tags linkfy-link
// @Accept json
// @Produce json
// @Param linkfy_id path string true "Linkfy profile ID"
// @Success 200 {object} DataResponse[[]LinkfyLink]
// @Failure 400 {object} MessageResponse
// @Failure 404 {object} MessageResponse
// @Failure 500 {object} MessageResponse
// @Router /linkfy/{linkfy_id}/links [get]
func (h *LinkfyLinkHandler) GetLinksByLinkfyID(c *gin.Context) {
	// Get linkfy_id from path
	linkfyIDStr := c.Param("linkfy_id")

	// Get links by linkfy ID
	links, custErr := h.repository.GetLinkByLinkfyID(linkfyIDStr)
	if custErr != nil {
		c.JSON(custErr.Code(), MessageResponse{
			Message: custErr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse[[](*LinkfyLink)]{
		Data:    links,
		Message: "Links retrieved successfully",
	})
}

package linkfy

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/yantology/linkfy/pkg/customerror"
)

type LinkfyService struct {
}

func NewLinkfyService() *LinkfyService {
	return &LinkfyService{}
}

func (s *LinkfyService) UsernameSanitize(username string) *customerror.CustomError {
	// Check for empty username
	if len(username) == 0 {
		return customerror.NewCustomError(nil, "Username cannot be empty", http.StatusBadRequest)
	}

	// Check for spaces
	if strings.Contains(username, " ") {
		return customerror.NewCustomError(nil, "Username cannot contain spaces", http.StatusBadRequest)
	}

	// Check if username contains only allowed characters
	matched, err := regexp.MatchString(`^[a-zA-Z0-9.-]+$`, username)
	if err != nil {
		return customerror.NewCustomError(err, "Error validating username", http.StatusInternalServerError)
	}
	if !matched {
		return customerror.NewCustomError(nil, "Username can only contain letters, numbers, dots, and hyphens", http.StatusBadRequest)
	}

	return nil
}

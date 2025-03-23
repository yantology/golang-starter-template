package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yantology/retail-pro-be/pkg/resendutils"
)

type authHandler struct {
	authService    AuthService
	authRepository *AuthRepository
	emailSender    resendutils.ResendUtilsInterface
	emailTemplate  EmailTemplateInterface
}

func NewAuthHandler(
	authService AuthService,
	authRepository *AuthRepository,
	emailSender resendutils.ResendUtilsInterface,
	emailTemplate EmailTemplateInterface,
) *authHandler {
	return &authHandler{
		authService:    authService,
		authRepository: authRepository,
		emailSender:    emailSender,
		emailTemplate:  emailTemplate,
	}
}

// @Summary Request activation token
// @Description Request a token for registration or password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param type path string true "Token type (registration or forget-password)"
// @Param request body TokenRequest true "Token request parameters"
// @Success 200 {object} MessageResponse "Success response with message"
// @Failure 400 {object} MessageResponse "Bad request response"
// @Failure 404 {object} MessageResponse "Not found response"
// @Failure 409 {object} MessageResponse "Conflict response"
// @Router /auth/token/{type} [post]
func (h *authHandler) RequestToken(c *gin.Context) {
	tokenType := c.Param("type")

	// Validate token type
	if tokenType != "registration" && tokenType != "forget-password" {
		c.JSON(http.StatusBadRequest, MessageResponse{

			Message: "Tipe token tidak valid",
		})
		return
	}

	var req TokenRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{

			Message: "Format request tidak valid",
		})
		return
	}

	// Validate email format
	if cuserr := h.authService.ValidateEmail(req.Email); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Check if email exists based on token type
	if tokenType == "registration" {
		if cuserr := h.authRepository.CheckExistingEmail(req.Email); cuserr != nil {
			c.JSON(http.StatusConflict, MessageResponse{

				Message: "test 1 2 3",
			})
			return
		}
	} else if tokenType == "forget-password" {
		if cuserr := h.authRepository.CheckExistingEmail(req.Email); cuserr == nil {
			c.JSON(http.StatusNotFound, MessageResponse{

				Message: "Email tidak terdaftar",
			})
			return
		}
	}

	// Generate activation token
	token, cuserr := h.authService.GenerateActivationToken()
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Hash the token before storing
	hashedToken, cuserr := h.authService.HashString(token)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Save token to database
	tokenReq := &ActivationTokenRequest{
		Email:          req.Email,
		ActivationCode: hashedToken,
		TokenType:      tokenType,
		ExpiryMinutes:  15, // 15 minutes expiry as per docs
	}

	if cuserr := h.authRepository.SaveActivationToken(tokenReq); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{
			Message: "Gagal menyimpan token",
		})
		return
	}

	// Generate email content based on token type
	var emailHTML, emailSubject string
	if tokenType == "registration" {
		emailHTML = h.emailTemplate.GenerateRegistrationEmail(req.Email, token)
		emailSubject = "Kode Aktivasi Pendaftaran"
	}
	if tokenType == "forget-password" {
		emailHTML = h.emailTemplate.GeneratePasswordResetEmail(req.Email, token)
		emailSubject = "Kode Reset Password"
	}

	// Send email
	if cuserr := h.emailSender.Send(emailHTML, emailSubject, []string{req.Email}); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{

		Message: "Kode aktivasi telah dikirim ke email",
	})
}

// @Summary Register new user
// @Description Register a new user with activation code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} MessageResponse "Success response with message"
// @Failure 400 {object} MessageResponse "Bad request response"
// @Failure 401 {object} MessageResponse "Unauthorized response"
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{

			Message: "Format request tidak valid",
		})
		return
	}

	// Validate registration input
	regReq := RegistrationRequest{
		Email:                req.Email,
		Username:             req.Fullname,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}

	if cuserr := h.authService.ValidateRegistrationInput(regReq); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Validate activation token
	tokenReq := &GetActivationTokenRequest{
		Email:     req.Email,
		TokenType: "registration",
	}

	token, cuserr := h.authRepository.GetActivationToken(tokenReq)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Verify token
	if cuserr := h.authService.VerifyHash(token, req.ActivationCode); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Hash password
	hashedPassword, cuserr := h.authService.HashString(req.Password)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Create user
	createUserReq := &CreateUserRequest{
		Email:        req.Email,
		Fullname:     req.Fullname,
		PasswordHash: hashedPassword,
	}

	if cuserr := h.authRepository.CreateUser(createUserReq); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{
		Message: "Pendaftaran berhasil, silakan login"})
}

// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} DataResponse[jwtResponseData] "Success response with JWT tokens"
// @Failure 400 {object} MessageResponse "Bad request response"
// @Failure 401 {object} MessageResponse "Unauthorized response"
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req LoginRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{

			Message: "Format request tidak valid",
		})
		return
	}

	// Get user by email
	user, cuserr := h.authRepository.GetUserByEmail(req.Email)
	if cuserr != nil {
		c.JSON(http.StatusUnauthorized, MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Verify password
	if cuserr := h.authService.VerifyHash(user.PasswordHash, req.Password); cuserr != nil {
		c.JSON(http.StatusUnauthorized, MessageResponse{

			Message: "hash tidak valid",
		})
		return
	}

	// Generate token pair
	tokenPairReq := TokenPairRequest{
		UserID: user.ID,
		Email:  user.Email,
	}

	tokenPair, cuserr := h.authService.GenerateTokenPair(tokenPairReq)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse[jwtResponseData]{
		Message: "Login berhasil",
		Data: jwtResponseData{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    900, // 15 minutes in seconds
		},
	})
}

// @Summary Reset password
// @Description Reset user password using activation code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgetPasswordRequest true "Password reset details"
// @Success 200 {object} MessageResponse "Success response with message"
// @Failure 400 {object} MessageResponse "Bad request response"
// @Failure 401 {object} MessageResponse "Unauthorized response"
// @Router /auth/forget-password [post]
func (h *authHandler) ForgetPassword(c *gin.Context) {
	var req ForgetPasswordRequest
	if cuserr := c.ShouldBindJSON(&req); cuserr != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{

			Message: "Format request tidak valid",
		})
		return
	}

	// Validate password match
	if cuserr := h.authService.ValidatePasswordInput(req.NewPassword, req.NewPasswordConfirmation); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Validate activation token
	tokenReq := &GetActivationTokenRequest{
		Email:     req.Email,
		TokenType: "forget-password",
	}

	token, cuserr := h.authRepository.GetActivationToken(tokenReq)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Verify token
	if cuserr := h.authService.VerifyHash(token, req.ActivationCode); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{
			Message: cuserr.Message(),
		})
		return
	}

	// Hash new password
	hashedPassword, cuserr := h.authService.HashString(req.NewPassword)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	// Update user password
	updateReq := &UpdatePasswordRequest{
		Email:           req.Email,
		NewPasswordHash: hashedPassword,
	}

	if cuserr := h.authRepository.UpdateUserPassword(updateReq); cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{

		Message: "Password berhasil diubah",
	})
}

// @Summary Refresh token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token details"
// @Success 200 {object} DataResponse[jwtResponseData] "Success response with new JWT tokens"
// @Failure 400 {object} MessageResponse "Bad request response"
// @Failure 401 {object} MessageResponse "Unauthorized response"
// @Router /auth/refresh-token [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	// Get refresh token from cookies
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "Refresh token tidak ditemukan dalam cookies",
		})
		return
	}

	// Create request object with token from cookie
	req := RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	// Validate refresh token
	claims, cuserr := h.authService.ValidateRefreshTokenClaims(req.RefreshToken)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	tokenPair := TokenPairRequest{
		UserID: claims.UserID,
		Email:  claims.Email,
	}

	// Generate new access token
	accessToken, cuserr := h.authService.GenerateTokenPair(tokenPair)
	if cuserr != nil {
		c.JSON(cuserr.Code(), MessageResponse{

			Message: cuserr.Message(),
		})
		return
	}

	c.JSON(http.StatusOK, DataResponse[jwtResponseData]{
		Message: "Token berhasil diperbarui",
		Data: jwtResponseData{
			AccessToken:  accessToken.AccessToken,
			RefreshToken: accessToken.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    900, // 15 minutes in seconds
		}})
}

// RegisterRoutes registers all auth routes
func (h *authHandler) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/token/:type", h.RequestToken)
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/forget-password", h.ForgetPassword)
		authGroup.GET("/refresh-token", h.RefreshToken)
	}
}

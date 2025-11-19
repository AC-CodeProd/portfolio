package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/config"
	"portfolio/domain"
	"portfolio/domain/usecases"
	authDto "portfolio/dto/auth"
	"portfolio/logger"
	"strings"
)

type authHandler struct {
	AbstractHandler
	authUseCase *usecases.AuthUseCase
	cfg         *config.JWTConfig
	logger      *logger.Logger
}

func NewAuthHandler(authUseCase *usecases.AuthUseCase, cfg *config.JWTConfig, logger *logger.Logger) []*routes.NamedRoute {
	authHandler := authHandler{
		authUseCase: authUseCase,
		cfg:         cfg,
		logger:      logger,
	}
	return []*routes.NamedRoute{
		{
			Name:    "LoginHandler",
			Pattern: "POST /auth/login",
			Handler: authHandler.Login,
		},
		{
			Name:    "LogoutHandler",
			Pattern: "DELETE /auth/logout",
			Handler: authHandler.Logout,
		},
	}
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate user and return JWT token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.AuthRequest									true	"Login request body"
//	@Success		200		{object}	shared.APIResponse{data=dto.AuthSuccess}		"Login successful"
//	@Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}	"Bad request"
//	@Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}	"Unauthorized"
//	@Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}	"Internal Server Error"
//	@Router			/admin/auth/login [post]
func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req authDto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("Failed to decode login request: %v", err)
		utils.WriteErrorResponse(w, domain.NewValidationError("body", "Invalid login request body", nil))
		return
	}
	resp, err := ah.authUseCase.Login(ctx, &req)
	if err != nil {
		ah.logger.Error("Login failed: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, resp)
}

// @Summary		User logout
// @Description	Logout user and invalidate JWT token
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		204	"No Content"
// @Failure		400		{object}	shared.APIResponse{errors=[]shared.APIError}	"Bad request"
// @Failure		401		{object}	shared.APIResponse{errors=[]shared.APIError}	"Unauthorized"
// @Failure		500		{object}	shared.APIResponse{errors=[]shared.APIError}	"Internal Server Error"
// @Router			/admin/auth/logout [delete]
func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenString := cleanTokenString(r.Header.Get("Authorization"))
	userID, ok := ah.getUserIDFromContext(w, r)
	if !ok {
		ah.logger.Error("Failed to get user ID from context")
		return
	}
	fmt.Println("Token:", tokenString, userID)
	err := ah.authUseCase.RevokeToken(ctx, userID, tokenString)
	if err != nil {
		ah.logger.Error("Failed to revoke token: %v", err)
		utils.WriteErrorResponse(w, err)
		return
	}
	utils.WriteSuccessResponse(w, http.StatusNoContent, nil)
}

func cleanTokenString(tokenString string) string {
	return strings.TrimSpace(strings.ReplaceAll(tokenString, "Bearer ", ""))
}

package middlewares

import (
	"context"
	"encoding/base64"
	"net/http"
	"portfolio/api/http/utils"
	"portfolio/config"
	"portfolio/domain"
	"portfolio/domain/entities"
	"portfolio/domain/usecases"
	"portfolio/logger"
	"portfolio/shared"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userCtxKey struct{}

var userKey = &userCtxKey{}

type AuthMiddleware struct {
	skipPaths   []string
	authUseCase *usecases.AuthUseCase
	cfg         *config.JWTConfig
	logger      *logger.Logger
}
type AuthMiddlewareOption func(*AuthMiddleware)

func NewAuthMiddleware(authUseCase *usecases.AuthUseCase, logger *logger.Logger, cfg *config.JWTConfig, options ...AuthMiddlewareOption) *AuthMiddleware {
	authMiddleware := &AuthMiddleware{
		authUseCase: authUseCase,
		cfg:         cfg,
		logger:      logger,
	}

	for _, option := range options {
		option(authMiddleware)
	}
	return authMiddleware
}

func AuthMiddlewareWithSkipPaths(skipPaths []string) AuthMiddlewareOption {
	return func(am *AuthMiddleware) {
		am.skipPaths = skipPaths
	}
}

func (am *AuthMiddleware) MiddlewareBasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.logger.Error("Authorization header is missing")
			am.writeUnauthorizedBasicAuth(w, domain.NewValidationError("Authorization header required", "authorization", nil))
			return
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			am.logger.Error("Invalid authorization format")
			am.writeUnauthorizedBasicAuth(w, domain.NewValidationError("Invalid authorization format", "authorization", nil))
			return
		}

		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			am.logger.Error("Failed to decode basic auth credentials")
			am.writeUnauthorizedBasicAuth(w, domain.NewValidationError("Invalid basic auth credentials", "authorization", nil))
			return
		}

		parts := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(parts) != 2 {
			am.logger.Error("Invalid basic auth credentials")
			am.writeUnauthorizedBasicAuth(w, domain.NewValidationError("Invalid basic auth credentials", "authorization", nil))
			return
		}

		username := parts[0]
		password := parts[1]

		if user, err := am.authUseCase.ValidateCredentials(r.Context(), username, password); err != nil {
			am.logger.Error("Failed to authenticate user")
			am.writeUnauthorizedBasicAuth(w, domain.NewValidationError("Invalid username or password", "authorization", nil))
			return
		} else {
			ctx := context.WithValue(r.Context(), userKey, user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) MiddlewareBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(am.skipPaths) > 0 {
			for _, path := range am.skipPaths {
				if strings.Contains(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.logger.Error("Authorization header is missing")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Authorization header required", "authorization", nil))
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			am.logger.Error("Invalid authorization format")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Invalid authorization format", "authorization", nil))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != am.cfg.SigningMethod {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(am.cfg.Secret), nil
		})
		if err != nil || !token.Valid {
			am.logger.Error("Invalid or expired token")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Invalid or expired token", "token", nil))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			am.logger.Error("Invalid token claims")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Invalid token claims", "token", nil))
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			am.logger.Error("Token expired")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Token expired", "token", nil))
			return
		}

		issuer, ok := claims["iss"].(string)
		if !ok || issuer != am.cfg.Issuer {
			am.logger.Error("Invalid token issuer")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Invalid token issuer", "token", nil))
			return
		}

		audience, ok := claims["aud"].(string)
		if !ok || audience != am.cfg.Audience {
			am.logger.Error("Invalid token audience")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Invalid token audience", "token", nil))
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			am.logger.Error("User ID not found in token")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("User ID not found in token", "token", nil))
			return
		}

		if revoked, err := am.authUseCase.IsTokenRevoked(r.Context(), int(userID), tokenString); err != nil {
			am.logger.Error("Failed to check token revocation")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Failed to check token revocation", "token", nil))
			return
		} else if revoked {
			am.logger.Error("Token has been revoked")
			am.writeUnauthorizedBearerToken(w, domain.NewValidationError("Token has been revoked", "token", nil))
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey{}, int(userID))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) writeUnauthorizedBasicAuth(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Basic realm="Documentation Area"`)
	w.WriteHeader(http.StatusUnauthorized)
}

func (am *AuthMiddleware) writeUnauthorizedBearerToken(w http.ResponseWriter, err error) {
	if domainErr, ok := domain.AsDomainError(err); ok {
		apiError := utils.DomainErrorToAPIError(domainErr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(domainErr.HTTPStatus())

		response := struct {
			Errors []*shared.APIError `json:"errors"`
		}{
			Errors: []*shared.APIError{apiError},
		}

		utils.JSONResponse(w, domainErr.HTTPStatus(), response)
		return
	}

	unauthorizedErr := domain.NewUnauthorizedError(err.Error())
	apiError := utils.DomainErrorToAPIError(unauthorizedErr)

	response := struct {
		Errors []*shared.APIError `json:"errors"`
	}{
		Errors: []*shared.APIError{apiError},
	}

	utils.JSONResponse(w, unauthorizedErr.HTTPStatus(), response)
}

func GetUserIDFromContext(r *http.Request) any {
	return r.Context().Value(userCtxKey{})
}

func GetUserFromContext(r *http.Request) (*entities.User, bool) {
	u, ok := r.Context().Value(userKey).(*entities.User)
	return u, ok
}

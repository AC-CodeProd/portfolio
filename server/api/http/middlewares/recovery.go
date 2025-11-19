package middlewares

import (
	"net/http"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/logger"
	"portfolio/shared"
	"runtime/debug"
)

func RecoveryMiddleware(logger *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered: %v", err)
					logger.Error("Stack trace: %s", string(debug.Stack()))

					domainErr := domain.NewInternalError("The server encountered an unexpected condition", nil)
					apiError := utils.DomainErrorToAPIError(domainErr)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(domainErr.HTTPStatus())

					response := struct {
						Errors []*shared.APIError `json:"errors"`
					}{
						Errors: []*shared.APIError{apiError},
					}

					utils.JSONResponse(w, domainErr.HTTPStatus(), response)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

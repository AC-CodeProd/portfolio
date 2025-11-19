package utils

import (
	"context"
	"net/http"
	"portfolio/domain"
	"portfolio/domain/usecases"
	"portfolio/shared"

	"github.com/google/uuid"
)

func GenerateRequestID() string {
	return uuid.New().String()
}

func SetRequestIDInContext(r *http.Request) (context.Context, string) {
	requestID := GenerateRequestID()
	ctx := context.WithValue(r.Context(), shared.REQUEST_ID_KEY, requestID)
	return ctx, requestID
}

func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(shared.REQUEST_ID_KEY).(string); ok {
		return requestID
	}
	return "unknown"
}

func GetPortfolioOwnerID(settingUseCase *usecases.SettingUseCase, ctx context.Context, w http.ResponseWriter) (int, error) {
	settings, err := settingUseCase.GetSettings(ctx)

	if err != nil {
		WriteErrorResponse(w, err)
		return 0, err
	}

	if settings == nil {
		WriteErrorResponse(w, domain.NewInternalError("Settings not initialized", nil))
		return 0, domain.NewInternalError("Settings not initialized", nil)
	}

	if settings.PortfolioOwnerID <= 0 {
		WriteErrorResponse(w, domain.NewInternalError("PortfolioOwnerID not configured in settings", nil))
		return 0, domain.NewInternalError("PortfolioOwnerID not configured in settings", nil)
	}

	return settings.PortfolioOwnerID, nil
}

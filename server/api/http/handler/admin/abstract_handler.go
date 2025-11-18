package admin

import (
	"net/http"
	"portfolio/api/http/middlewares"
	"portfolio/api/http/utils"
	"portfolio/domain"
	"portfolio/domain/usecases"
)

type AbstractHandler struct {
	settingUseCase *usecases.SettingUseCase
}

func (ah *AbstractHandler) getUserIDFromContext(w http.ResponseWriter, r *http.Request) (int, bool) {
	userIDFromContext := middlewares.GetUserIDFromContext(r)
	userID, ok := userIDFromContext.(int)
	if !ok || userID <= 0 {
		utils.WriteErrorResponse(w, domain.NewValidationError("Invalid user ID in context", "user_id", nil))
		return 0, false
	}
	return userID, true
}

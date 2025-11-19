package handler

import (
	"portfolio/domain/usecases"
)

type AbstractHandler struct {
	settingUseCase *usecases.SettingUseCase
}

// func (ah *AbstractHandler) getUserFromContext(w http.ResponseWriter, r *http.Request) (*entities.User, bool) {
// 	return middlewares.GetUserFromContext(r)
// }

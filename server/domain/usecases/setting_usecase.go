package usecases

import (
	"context"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type SettingUseCase struct {
	settingRepo interfaces.SettingRepository
	logger      *logger.Logger
}

func NewSettingUseCase(settingRepo interfaces.SettingRepository, logger *logger.Logger) *SettingUseCase {
	return &SettingUseCase{
		settingRepo: settingRepo,
		logger:      logger,
	}
}

func (suc *SettingUseCase) Upsert(ctx context.Context, settingJson *entities.SettingJson) error {
	return suc.settingRepo.Upsert(ctx, settingJson)
}

func (suc *SettingUseCase) GetSettings(ctx context.Context) (*entities.SettingJson, error) {
	return suc.settingRepo.GetSettings(ctx)
}

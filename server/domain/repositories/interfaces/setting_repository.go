package interfaces

import (
	"context"
	"portfolio/domain/entities"
)

type SettingRepository interface {
	GetSettings(ctx context.Context) (*entities.SettingJson, error)
	Upsert(ctx context.Context, settingJson *entities.SettingJson) error
	Delete(ctx context.Context) error
}

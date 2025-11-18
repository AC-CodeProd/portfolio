package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"portfolio/domain/entities"
	"portfolio/domain/repositories/interfaces"
	"portfolio/logger"
)

type settingRepository struct {
	settingKey string
	db         *sql.DB
	logger     *logger.Logger
}

func NewSettingRepository(db *sql.DB, logger *logger.Logger, settingKey string) interfaces.SettingRepository {
	return &settingRepository{
		db:         db,
		logger:     logger,
		settingKey: settingKey,
	}
}

func (repo *settingRepository) Upsert(ctx context.Context, settingJson *entities.SettingJson) error {

	dataBytes, err := json.Marshal(settingJson)

	if err != nil {
		return err
	}

	query := `
	INSERT INTO settings (setting_key, setting_json, setting_created_at, setting_updated_at)
	VALUES (?, JSON(?), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(setting_key) DO UPDATE SET
		setting_json = JSONB(excluded.setting_json),
		setting_updated_at = CURRENT_TIMESTAMP
	`
	_, err = repo.db.ExecContext(
		ctx, query,
		repo.settingKey,
		dataBytes,
	)
	return err
}

func (sr *settingRepository) GetSettings(ctx context.Context) (*entities.SettingJson, error) {
	query := `SELECT setting_key, JSON(setting_json), setting_created_at, setting_updated_at FROM settings WHERE setting_key = ?`
	var setting entities.Setting
	err := sr.db.QueryRowContext(ctx, query, sr.settingKey).Scan(
		&setting.SettingKey,
		&setting.SettingJson,
		&setting.SettingCreatedAt,
		&setting.SettingUpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var settingJson entities.SettingJson
	err = json.Unmarshal(setting.SettingJson, &settingJson)
	if err != nil {
		return nil, err
	}

	return &settingJson, nil
}

func (sr *settingRepository) Delete(ctx context.Context) error {
	query := `DELETE FROM settings WHERE setting_key = ?`
	_, err := sr.db.ExecContext(ctx, query, sr.settingKey)
	return err
}

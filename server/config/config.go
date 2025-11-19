package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"portfolio/helpers"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	JWT        JWTConfig      `yaml:"jwt"`
	Database   DatabaseConfig `yaml:"database"`
	CORS       CORSConfig     `yaml:"cors"`
	Logging    LoggingConfig  `yaml:"logging"`
	Server     ServerConfig   `yaml:"server"`
	Admin      AdminConfig    `yaml:"admin"`
	SettingKey string         `yaml:"setting_key"`
}

type ServerConfig struct {
	Port        string `yaml:"port"`
	Environment string `yaml:"environment"`
	Mode        string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver             string            `yaml:"driver"`
	Path               string            `yaml:"path"`
	MaxOpenConnections int               `yaml:"max_open_connections"`
	MaxIdleConnections int               `yaml:"max_idle_connections"`
	ConnMaxLifetime    int               `yaml:"conn_max_lifetime"`
	Options            map[string]string `yaml:"options"`
	Pragmas            map[string]string `yaml:"pragmas"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type LoggingConfig struct {
	File        string  `yaml:"file"`
	Level       string  `yaml:"level"`
	MaxSize     float32 `yaml:"max_size"`
	MaxBackups  int     `yaml:"max_backups"`
	MaxAge      int     `yaml:"max_age"`
	Compress    *bool   `yaml:"compress"`
	RotateDaily *bool   `yaml:"rotate_daily"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Salt     string `yaml:"salt"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	Expiration    string `yaml:"expiration"`
	Issuer        string `yaml:"issuer"`
	Audience      string `yaml:"audience"`
	SigningMethod string `yaml:"signing_method"`
}

func LoadConfig(configPath string) (*Config, error) {

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config := getDefaultConfig()
			saveDefaultConfig(configPath, config)
			return config, nil
		}
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	overrideWithEnv(&config)

	return &config, nil
}

func getDefaultConfig() *Config {
	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = "."
	}
	return &Config{
		Server: ServerConfig{
			Port: "3000",
			Mode: "development",
		},
		Database: DatabaseConfig{
			Driver:             "sqlite3",
			Path:               filepath.Join(baseDir, "data", "portfolio.sqlite3"),
			MaxOpenConnections: 25,
			MaxIdleConnections: 5,
			ConnMaxLifetime:    300,
			Pragmas: map[string]string{
				"foreign_keys":       "ON",
				"journal_mode":       "WAL",
				"synchronous":        "NORMAL",
				"busy_timeout":       "5000",
				"temp_store":         "MEMORY",
				"mmap_size":          "134217728", // 128MB
				"journal_size_limit": "67108864",  // 64MB
				"cache_size":         "2000",
			},
			Options: map[string]string{
				"_transaction_mode": "IMMEDIATE",
			},
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:5173",
				"http://localhost:4200",
			},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		Logging: LoggingConfig{
			File:        "logs/portfolio.log",
			Level:       "info",
			MaxSize:     10,
			MaxBackups:  7,
			MaxAge:      30,
			Compress:    helpers.BoolPtr(true),
			RotateDaily: helpers.BoolPtr(true),
		},
		Admin: AdminConfig{
			Username: "admin",
		},
		JWT: JWTConfig{
			Secret:        "your_jwt_secret_key",
			Expiration:    "24h",
			Issuer:        "portfolio-api",
			Audience:      "portfolio-client",
			SigningMethod: "HS256",
		},
		SettingKey: "portfolio",
	}
}

func saveDefaultConfig(path string, config *Config) {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling default config: %v", err))
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing default config file: %v", err))
	}

	fmt.Printf("Default config file created at: %s\n", path)
}

func overrideWithEnv(config *Config) {
	if envFile := os.Getenv("PORTFOLIO_ENV_FILE"); envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("Failed to load environment file: %v", err)
		}
	}

	if port := os.Getenv("PORTFOLIO_PORT"); port != "" {
		config.Server.Port = port
	}
	if environment := os.Getenv("PORTFOLIO_ENVIRONMENT"); environment != "" {
		config.Server.Environment = environment
	}
	if mode := os.Getenv("PORTFOLIO_MODE"); mode != "" {
		config.Server.Mode = mode
	}
	if logFile := os.Getenv("PORTFOLIO_LOG_FILE"); logFile != "" {
		config.Logging.File = logFile
	}
	if logLevel := os.Getenv("PORTFOLIO_LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	if jwtSecret := os.Getenv("PORTFOLIO_JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}
	if jwtExpiration := os.Getenv("PORTFOLIO_JWT_EXPIRATION"); jwtExpiration != "" {
		config.JWT.Expiration = jwtExpiration
	}
	if jwtIssuer := os.Getenv("PORTFOLIO_JWT_ISSUER"); jwtIssuer != "" {
		config.JWT.Issuer = jwtIssuer
	}
	if jwtAudience := os.Getenv("PORTFOLIO_JWT_AUDIENCE"); jwtAudience != "" {
		config.JWT.Audience = jwtAudience
	}

	if jwtSigningMethod := os.Getenv("PORTFOLIO_JWT_SIGNING_METHOD"); jwtSigningMethod != "" {
		config.JWT.SigningMethod = jwtSigningMethod
	}
	if salt := os.Getenv("PORTFOLIO_ADMIN_SALT"); salt != "" {
		config.Admin.Salt = salt
	}
	if settingKey := os.Getenv("PORTFOLIO_SETTING_KEY"); settingKey != "" {
		config.SettingKey = settingKey
	}
}

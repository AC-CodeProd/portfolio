/*
@title           My Personal Portfolio API
@version         0.0.1
@description		 This is the API documentation for my personal portfolio backend server.
@description.markdown
@contact.name	Alain CAJUSTE
@contact.url     https://portfolio.example.com
@contact.email   cajuste.alain@gmail.com
@host           localhost:3000
@accept json
@produce json
@schemes https
@basePath       /v1
@securityDefinitions.apikey BearerAuth
@in header
@name Authorization
@description Bearer {token}
*/
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"portfolio/api/http/handler"
	"portfolio/api/http/handler/admin"
	"portfolio/api/http/handler/doc"
	"portfolio/api/http/middlewares"
	"portfolio/api/http/routes"
	"portfolio/api/http/utils"
	"portfolio/config"
	"portfolio/domain/repositories/interfaces"
	"portfolio/domain/usecases"
	"portfolio/infrastructure/sqlite"
	"portfolio/logger"
	"portfolio/service"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type RepositoryBundle struct {
	Setting      interfaces.SettingRepository
	PersonalInfo interfaces.PersonalInfoRepository
	RevokeToken  interfaces.RevokedTokenRepository
	User         interfaces.UserRepository
	Project      interfaces.ProjectRepository
	Skill        interfaces.SkillRepository
	Experience   interfaces.ExperienceRepository
	Education    interfaces.EducationRepository
	Technology   interfaces.TechnologyRepository
}

type UseCaseBundle struct {
	Setting      *usecases.SettingUseCase
	PersonalInfo *usecases.PersonalInfoUseCase
	Auth         *usecases.AuthUseCase
	Project      *usecases.ProjectUseCase
	Skill        *usecases.SkillUseCase
	Experience   *usecases.ExperienceUseCase
	Education    *usecases.EducationUseCase
	Technology   *usecases.TechnologyUseCase
}

func initializeConfig() (*config.Config, *logger.Logger, error) {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		return nil, nil, err
	}

	logger := logger.NewLogger(&cfg.Logging, cfg.Server.Environment)
	return cfg, logger, nil
}

func initializeDatabase(cfg *config.Config, logger *logger.Logger) (*sql.DB, error) {
	logger.Info("Initializing database connection...")

	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, err
	}

	db, err := sqlite.NewConnection(&cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initializeRepositories(db *sql.DB, cfg *config.Config, logger *logger.Logger) *RepositoryBundle {
	logger.Info("Initializing repositories...")

	return &RepositoryBundle{
		Setting:      sqlite.NewSettingRepository(db, logger, cfg.SettingKey),
		PersonalInfo: sqlite.NewPersonalInfoRepository(db, logger),
		RevokeToken:  sqlite.NewRevokedTokenRepository(db, logger),
		User:         sqlite.NewUserRepository(db, logger),
		Project:      sqlite.NewProjectRepository(db, logger),
		Skill:        sqlite.NewSkillRepository(db, logger),
		Experience:   sqlite.NewExperienceRepository(db, logger),
		Education:    sqlite.NewEducationRepository(db, logger),
		Technology:   sqlite.NewTechnologyRepository(db, logger),
	}
}

func initializeUseCases(repos *RepositoryBundle, cfg *config.Config, logger *logger.Logger) *UseCaseBundle {
	logger.Info("Initializing use cases...")

	authService := service.NewAuthService(&cfg.JWT)
	settingUseCase := usecases.NewSettingUseCase(repos.Setting, logger)

	return &UseCaseBundle{
		Setting:      settingUseCase,
		PersonalInfo: usecases.NewPersonalInfoUseCase(repos.PersonalInfo, logger),
		Auth:         usecases.NewAuthUseCase(repos.User, repos.RevokeToken, settingUseCase, authService, logger, cfg.Admin.Salt),
		Project:      usecases.NewProjectUseCase(repos.Project, repos.User, repos.Setting, logger),
		Skill:        usecases.NewSkillUseCase(repos.Skill, repos.User, logger),
		Experience:   usecases.NewExperienceUseCase(repos.Experience, repos.User, logger),
		Education:    usecases.NewEducationUseCase(repos.Education, repos.User, logger),
		Technology:   usecases.NewTechnologyUseCase(repos.Technology, repos.User, logger),
	}
}

func setupMiddlewares(authUseCase *usecases.AuthUseCase, jwtConfig *config.JWTConfig, cfg *config.Config, logger *logger.Logger) (
	*middlewares.AuthMiddleware, *middlewares.RateLimiter, func(http.Handler) http.Handler,
	func(http.Handler) http.Handler, func(http.Handler) http.Handler, func(http.Handler) http.Handler) {

	authMiddleware := middlewares.NewAuthMiddleware(authUseCase, logger, jwtConfig, middlewares.AuthMiddlewareWithSkipPaths(
		[]string{"/auth/login"},
	))
	rateLimiter := middlewares.NewRateLimiter(time.Second, 10)
	if cfg.Logging.Level == "debug" {
		rateLimiter.SetLogger(logger)
	}

	loggingMW := middlewares.LoggingMiddleware(logger)
	corsMW := middlewares.CORSMiddleware(cfg)
	recoveryMW := middlewares.RecoveryMiddleware(logger)
	responseMW := middlewares.ResponseMiddleware(logger)

	return authMiddleware, rateLimiter, loggingMW, corsMW, recoveryMW, responseMW
}

func setupHandlers(
	settingUseCase *usecases.SettingUseCase,
	personalInfoUseCase *usecases.PersonalInfoUseCase,
	authUseCase *usecases.AuthUseCase,
	projectUseCase *usecases.ProjectUseCase,
	skillUseCase *usecases.SkillUseCase,
	experienceUseCase *usecases.ExperienceUseCase,
	educationUseCase *usecases.EducationUseCase,
	technologyUseCase *usecases.TechnologyUseCase,
	jwtConfig *config.JWTConfig,
	logger *logger.Logger,
) ([]*routes.NamedRoute, []*routes.NamedRoute) {

	personalInfoHandler := handler.NewPersonalInfoHandler(settingUseCase, personalInfoUseCase, logger)
	projectHandler := handler.NewProjectHandler(settingUseCase, projectUseCase, logger)
	skillHandler := handler.NewSkillHandler(settingUseCase, skillUseCase, logger)
	experienceHandler := handler.NewExperienceHandler(settingUseCase, experienceUseCase, logger)
	educationHandler := handler.NewEducationHandler(settingUseCase, educationUseCase, logger)
	technologyHandler := handler.NewTechnologyHandler(settingUseCase, technologyUseCase, logger)
	settingHandler := handler.NewSettingHandler(settingUseCase, logger)

	adminAuthHandler := admin.NewAuthHandler(authUseCase, jwtConfig, logger)
	adminPersonalInfoHandler := admin.NewPersonalInfoHandler(settingUseCase, personalInfoUseCase, logger)
	adminProjectHandler := admin.NewProjectHandler(settingUseCase, projectUseCase, logger)
	adminSkillHandler := admin.NewSkillHandler(settingUseCase, skillUseCase, logger)
	adminExperienceHandler := admin.NewExperienceHandler(settingUseCase, experienceUseCase, logger)
	adminEducationHandler := admin.NewEducationHandler(settingUseCase, educationUseCase, logger)
	adminTechnologyHandler := admin.NewTechnologyHandler(settingUseCase, technologyUseCase, logger)
	adminSettingHandler := admin.NewSettingHandler(settingUseCase, logger)

	var allAdminRoutes []*routes.NamedRoute
	allAdminRoutes = append(allAdminRoutes, adminAuthHandler...)
	allAdminRoutes = append(allAdminRoutes, adminPersonalInfoHandler...)
	allAdminRoutes = append(allAdminRoutes, adminProjectHandler...)
	allAdminRoutes = append(allAdminRoutes, adminSkillHandler...)
	allAdminRoutes = append(allAdminRoutes, adminExperienceHandler...)
	allAdminRoutes = append(allAdminRoutes, adminEducationHandler...)
	allAdminRoutes = append(allAdminRoutes, adminTechnologyHandler...)
	allAdminRoutes = append(allAdminRoutes, adminSettingHandler...)

	var allRoutes []*routes.NamedRoute
	allRoutes = append(allRoutes, personalInfoHandler...)
	allRoutes = append(allRoutes, projectHandler...)
	allRoutes = append(allRoutes, skillHandler...)
	allRoutes = append(allRoutes, experienceHandler...)
	allRoutes = append(allRoutes, educationHandler...)
	allRoutes = append(allRoutes, technologyHandler...)
	allRoutes = append(allRoutes, settingHandler...)

	return allRoutes, allAdminRoutes
}

func setupHTTPServer(useCases *UseCaseBundle, cfg *config.Config, logger *logger.Logger) *http.Server {
	logger.Info("Setting up HTTP server...")

	authMiddleware, rateLimiter, loggingMW, corsMW, recoveryMW, responseMW := setupMiddlewares(useCases.Auth, &cfg.JWT, cfg, logger)

	allRoutes, allAdminRoutes := setupHandlers(
		useCases.Setting,
		useCases.PersonalInfo, useCases.Auth, useCases.Project, useCases.Skill,
		useCases.Experience, useCases.Education, useCases.Technology, &cfg.JWT, logger,
	)
	docs := doc.NewDocsHandler(logger)

	baseMux := routes.SetupRoutes(allRoutes...)
	adminMux := routes.SetupRoutes(allAdminRoutes...)
	docsMux := routes.SetupRoutes(docs...)

	mux := http.NewServeMux()

	baseChain := middlewares.ChainMiddleware(
		responseMW,
		recoveryMW,
		corsMW,
		rateLimiter.Middleware,
		loggingMW,
	)

	adminChain := middlewares.ChainMiddleware(
		baseChain,
		authMiddleware.MiddlewareBearerToken,
	)

	docsChain := middlewares.ChainMiddleware(
		authMiddleware.MiddlewareBasicAuth,
	)

	mux.Handle("/v1/", baseChain(http.StripPrefix("/v1", baseMux)))
	mux.Handle("/admin/", adminChain(http.StripPrefix("/admin", adminMux)))
	mux.Handle("/doc/", docsChain(http.StripPrefix("/doc", docsMux)))

	mux.Handle("GET /health", baseChain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.JSONResponse(w, http.StatusOK, map[string]any{
				"message":   "OK",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}),
	))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	return server
}

func startServer(server *http.Server, useCases *UseCaseBundle, cfg *config.Config, logger *logger.Logger) error {
	logger.Info("Starting portfolio backend server...")

	ctx := context.Background()
	err := useCases.Auth.CreateDefaultAdmin(ctx, cfg.Admin.Username)
	if err != nil {
		logger.Error("Failed to create default admin: %v", err)
	}

	go func() {
		logger.Info("=== Portfolio Backend Server ===")
		logger.Info("Port: %s", cfg.Server.Port)
		logger.Info("Environment: %s", cfg.Server.Environment)
		logger.Info("Mode: %s", cfg.Server.Mode)
		logger.Info("Database: %s", cfg.Database.Path)
		logger.Info("===============================")
		logger.Info("🚀 API: http://localhost:%s/v1/", cfg.Server.Port)
		logger.Info("👑 Admin: http://localhost:%s/admin/", cfg.Server.Port)
		logger.Info("📚 Documentation: http://localhost:%s/doc/", cfg.Server.Port)
		logger.Info("💖 Health: http://localhost:%s/health", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
		return err
	} else {
		logger.Info("Server exited gracefully")
	}

	return nil
}

func main() {
	_ = godotenv.Load()

	cfg, logger, err := initializeConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	logger.Info("Starting portfolio backend server...")

	db, err := initializeDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Printf("db close: %v", err)
		}
	}()

	repos := initializeRepositories(db, cfg, logger)

	useCases := initializeUseCases(repos, cfg, logger)

	server := setupHTTPServer(useCases, cfg, logger)

	if err := startServer(server, useCases, cfg, logger); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}

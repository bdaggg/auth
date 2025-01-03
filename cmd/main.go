package main

import (
	"log"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/oauth"
	"auth-service/internal/handlers"
	"auth-service/internal/infrastructure/cache"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/middleware"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/security"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	// Database bağlantısı
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kurulamadı: %v", err)
	}
	defer db.Close()

	// Redis bağlantısı
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Redis bağlantısı kurulamadı: %v", err)
	}
	defer redisClient.Close()

	// JWT manager
	jwtManager := security.NewJWTManager(security.JWTConfig{
		AccessTokenSecret:  cfg.JWT.AccessTokenSecret,
		RefreshTokenSecret: cfg.JWT.RefreshTokenSecret,
		AccessTokenTTL:     cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL:    cfg.JWT.RefreshTokenTTL,
		Issuer:             cfg.JWT.Issuer,
	})

	// OAuth2 providers
	googleProvider := oauth.NewGoogleProvider(
		cfg.OAuth.GoogleClientID,
		cfg.OAuth.GoogleClientSecret,
		cfg.OAuth.GoogleRedirectURL,
	)

	// Email service
	emailService := service.NewEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	// TOTP service
	totpService := service.NewTOTPService(cfg.JWT.Issuer)

	// Repositories
	userRepo := repository.NewGormUserRepository(db.GetDB())
	sessionRepo := repository.NewSessionRepository(redisClient.GetClient())
	auditRepo := repository.NewAuditRepository(db.GetDB())
	securityRepo := repository.NewSecurityRepository(db.GetDB())

	// Services
	monitoringService := service.NewMonitoringService(auditRepo, securityRepo)
	securityService := service.NewSecurityService(
		userRepo,
		securityRepo,
		sessionRepo,
		monitoringService,
	)

	authService := service.NewAuthService(
		userRepo,
		jwtManager,
		emailService,
		totpService,
		sessionRepo,
		auditRepo,
		securityService,
		googleProvider,
	)

	// Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware'ler
	app.Use(middleware.RateLimit(redisClient.GetClient(), 100, time.Minute))
	app.Use(middleware.RequestMetrics(monitoringService))
	app.Use(recover.New())
	app.Use(logger.New())

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public routes
	auth := v1.Group("/auth")
	auth.Post("/register", handlers.Register(authService))
	auth.Post("/login", handlers.Login(authService))
	auth.Post("/refresh", handlers.RefreshToken(authService))
	auth.Post("/forgot-password", handlers.ForgotPassword(authService))
	auth.Post("/reset-password", handlers.ResetPassword(authService))
	auth.Get("/verify-email", handlers.VerifyEmail(authService))

	// OAuth routes
	auth.Get("/google/login", handlers.GoogleLogin(googleProvider))
	auth.Get("/google/callback", handlers.GoogleCallback(authService))

	// Protected routes
	protected := v1.Group("/protected")
	protected.Use(middleware.JWTAuth(jwtManager))

	// User routes
	user := protected.Group("/user")
	user.Post("/change-password", handlers.ChangePassword(authService))
	user.Post("/2fa/enable", handlers.Enable2FA(authService))
	user.Post("/2fa/verify", handlers.Verify2FA(authService))
	user.Get("/audit-logs", handlers.GetAuditLogs(authService))

	// Security routes
	security := protected.Group("/security")
	security.Use(middleware.RequireRole(entity.RoleSecurityAdmin))
	security.Post("/users/:id/block", handlers.BlockUser(securityService))
	security.Post("/users/:id/unblock", handlers.UnblockUser(securityService))
	security.Get("/alerts", handlers.GetSecurityAlerts(securityService))
	security.Get("/suspicious", handlers.GetSuspiciousActivities(securityService))

	// Monitoring routes
	monitoring := protected.Group("/monitoring")
	monitoring.Use(middleware.RequireRole(entity.RoleSystemMonitor))
	monitoring.Get("/metrics", handlers.GetMetrics(monitoringService))
	monitoring.Get("/active-users", handlers.GetActiveUsers(monitoringService))
	monitoring.Get("/blocked-users", handlers.GetBlockedUsers(monitoringService))

	// Admin routes
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(entity.RoleAdmin))
	admin.Get("/users", handlers.ListUsers(authService))
	admin.Post("/users/:id/role", handlers.ChangeUserRole(authService))

	log.Fatal(app.Listen(cfg.Server.Address))
}

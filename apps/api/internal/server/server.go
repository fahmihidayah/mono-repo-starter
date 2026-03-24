package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	"github.com/fahmihidayah/go-api-orchestrator/internal/controller"
	"github.com/fahmihidayah/go-api-orchestrator/internal/db"
	"github.com/fahmihidayah/go-api-orchestrator/internal/mail"
	"github.com/fahmihidayah/go-api-orchestrator/internal/middleware"
	"github.com/fahmihidayah/go-api-orchestrator/internal/repository"
	"github.com/fahmihidayah/go-api-orchestrator/internal/service"
	"github.com/fahmihidayah/go-api-orchestrator/internal/storage"
)

type Server struct {
	port               int
	config             *config.Config
	authMiddleware     func(http.Handler) http.Handler
	authController     *controller.AuthUserController
	userController     *controller.UserController
	postController     *controller.PostController
	categoryController *controller.CategoryController
	mediaController    *controller.MediaController
	roleController     *controller.RoleController
	// React Admin controllers
}

func NewServer() *http.Server {
	config, _ := config.LoadConfig()
	db, _ := db.NewDatabase(config)
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	mailer := mail.MailerProvider(config, mail.TemplateEngineProvider())

	userRepository := repository.UserRepositoryProvider(db)
	roleRepository := repository.RoleRepositoryProvider(db)
	tokenBlacklistRepository := repository.TokenBlacklistRepositoryProvider(db)
	postRepository := repository.PostRepositoryProvider(db)
	categoryRepository := repository.CategoryRepositoryProvider(db)
	mediaRepository := repository.MediaRepositoryProvider(db)

	userService := service.UserServiceProvider(
		userRepository, tokenBlacklistRepository, mailer, config,
	)

	postService := service.PostServiceProvider(
		postRepository, categoryRepository,
	)

	categoryService := service.CategoryServiceProvider(
		categoryRepository, config,
	)

	roleService := service.RoleServiceProvider(
		roleRepository, config,
	)

	// Initialize storage (local or S3 based on config)
	storageInstance, err := storage.StorageProvider(config)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	mediaService := service.MediaServiceProvider(
		mediaRepository, config, storageInstance,
	)

	NewServer := &Server{
		port:           port,
		config:         config,
		authMiddleware: middleware.AuthMiddleware(config, tokenBlacklistRepository),
		authController: controller.AuthUserControllerProvider(
			userService,
		),
		userController: controller.UserControllerProvider(
			userService,
		),
		postController: controller.PostControllerProvider(
			postService,
		),
		categoryController: controller.CategoryControllerProvider(
			categoryService,
		),
		mediaController: controller.MediaControllerProvider(
			mediaService,
		),
		roleController: controller.RoleControllerProvider(
			roleService,
		),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

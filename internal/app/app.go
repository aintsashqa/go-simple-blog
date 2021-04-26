package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aintsashqa/go-simple-blog/internal/config"
	"github.com/aintsashqa/go-simple-blog/internal/delivery/http"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/internal/serializer"
	"github.com/aintsashqa/go-simple-blog/internal/server"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/aintsashqa/go-simple-blog/internal/store"
	"github.com/aintsashqa/go-simple-blog/pkg/auth/jwt"
	"github.com/aintsashqa/go-simple-blog/pkg/cache/redis"
	"github.com/aintsashqa/go-simple-blog/pkg/database/mysql"
	"github.com/aintsashqa/go-simple-blog/pkg/hash/bcrypt"
	standart "github.com/aintsashqa/go-simple-blog/pkg/logger/standard"
	"github.com/aintsashqa/go-simple-blog/seeds"
)

// @title Go Simple Blog API
// @version 1.0.0
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func Run() {
	ctx := context.Background()

	logger := standart.NewStandartLoggerProvider()
	logger.Info("Initialize config")
	cfg, err := config.Init("config")
	if err != nil {
		logger.Critical(err)
	}

	logger.Info("Initialize database connection")
	database, err := mysql.NewMySQLProvider(mysql.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Username: cfg.Database.Username,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DatabaseName,
		Charset:  cfg.Database.Charset,
	})
	if err != nil {
		logger.Critical(err)
	}

	if cfg.Development {
		logger.Info("Seeding database")
		seeder := seeds.NewSeeder(database)
		if err := seeder.Seed(ctx, seeds.UserSeed, seeds.PostSeed); err != nil {
			logger.Critical(err)
		}
	}

	logger.Info("Initialize cache")
	cache, err := redis.NewRedisProvider(ctx, redis.Config{
		Host:     cfg.Cache.Host,
		Port:     cfg.Cache.Port,
		Username: cfg.Cache.Username,
		Password: cfg.Cache.Password,
		Database: cfg.Cache.Database,
	}, cfg.Cache.Expires)
	if err != nil {
		logger.Critical(err)
	}

	logger.Info("Initialize dependecies")
	repos := repository.NewRepository(database)
	serializer := serializer.NewSerializer()
	store := store.NewCacheStore(repos, cache, serializer)
	hasher := bcrypt.NewBcryptProvider()
	auth := jwt.NewJWTAuthorizationProvider(cfg.Auth.JWTSigningKey)

	services := service.NewService(service.ServiceDependencies{
		Logger:                        logger,
		DataProvider:                  store,
		Hasher:                        hasher,
		Authorization:                 auth,
		AuthorizationTokenExpiresTime: cfg.Auth.JWTExpiresTime,
	})

	handler := http.NewHandler(services)

	logger.Info("Starting server")
	srv := server.NewServer(cfg.App, handler.Init(cfg.App.Host, cfg.App.Port))
	go func() {
		if err := srv.Run(); err != nil {
			logger.Critical(err)
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := database.Close(); err != nil {
		logger.Critical(err)
	}
}

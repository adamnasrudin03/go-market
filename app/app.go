package app

import (
	"github.com/adamnasrudin03/go-market/app/controller"
	"github.com/adamnasrudin03/go-market/app/repository"
	"github.com/adamnasrudin03/go-market/app/service"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/adamnasrudin03/go-market/pkg/driver"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func WiringRepository(db *gorm.DB, cache *driver.RedisClient, cfg *configs.Configs, logger *logrus.Logger) *repository.Repositories {
	return &repository.Repositories{
		Auth:  repository.NewAuthRepository(db, cfg, logger),
		Cache: repository.NewCacheRepository(*cache, cfg, logger),
	}
}

func WiringService(repo *repository.Repositories, cfg *configs.Configs, logger *logrus.Logger) *service.Services {
	return &service.Services{
		Auth: service.NewAuthService(service.AuthSrv{Repo: repo.Auth, Cfg: cfg, Logger: logger}),
	}
}

func WiringController(srv *service.Services, cfg *configs.Configs, logger *logrus.Logger, validator *validator.Validate) *controller.Controllers {
	return &controller.Controllers{
		Auth: controller.NewAuthController(srv.Auth, logger, validator),
	}
}

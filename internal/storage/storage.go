package storage

import (
	"go.uber.org/zap"
	"just_for_fun/internal/config"
	"just_for_fun/internal/storage/db"
	"just_for_fun/pkg"
	"just_for_fun/pkg/logging"
)

type Storage struct {
	DB  *db.HelperDB
	Log *zap.Logger
}

var ModuleName = pkg.GetPackageName()

func New(cfg *config.Config, logger *logging.DynamicLogger) *Storage {
	logger.AddModule(ModuleName)
	log := logger.GetLogger(ModuleName)

	return &Storage{
		DB:  db.New(cfg, *log),
		Log: log,
	}
}

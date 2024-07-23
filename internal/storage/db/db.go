package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"just_for_fun/internal/config"
	"just_for_fun/internal/storage/db/repo"
)

type HelperDB struct {
	DB   *gorm.DB
	log  *zap.Logger
	User *repo.UserRepo
}

func New(cfg *config.Config, log zap.Logger) *HelperDB {
	dbConn := dbConnect(&log, dbDialect(cfg))
	userRepo := repo.NewUserRepo(log, dbConn)
	helper := &HelperDB{
		DB:   dbConn,
		log:  &log,
		User: userRepo,
	}

	log.Info("connected to database...")
	helper.setupStoragePool(cfg)
	return helper
}

func (hdb *HelperDB) setupStoragePool(cfg *config.Config) {
	sqlDB, err := hdb.DB.DB()
	if err != nil {
		hdb.log.Fatal("failed to connect database", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(cfg.DB.Pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.Pool.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(cfg.DB.Pool.IdleTimeout)
}

func dbConnect(log *zap.Logger, dbDialect gorm.Dialector) *gorm.DB {
	dbConn, err := gorm.Open(dbDialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}
	return dbConn
}

func dbDialect(cfg *config.Config) gorm.Dialector {
	return postgres.Open(cfg.DB.ConnString())
}

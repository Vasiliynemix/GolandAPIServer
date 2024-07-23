package main

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"just_for_fun/internal/config"
	"just_for_fun/internal/server"
	"just_for_fun/internal/storage"
	"just_for_fun/pkg"
	"just_for_fun/pkg/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	levelsUpToRootDir = 3
)

var moduleName = pkg.GetPackageName()

// @title Just for fun API
// @version 1.0
// @description Just for fun API description

// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.Load(levelsUpToRootDir)

	logger := logging.InitLogger(cfg)
	defer logger.Sync()

	logger.AddModule(moduleName)
	log := logger.GetLogger(moduleName)

	log.Info("Starting...", zap.Any("config", cfg))

	go checkLogFiles(logger)

	router := chi.NewRouter()

	storageHelper := storage.New(cfg, logger)

	serv := server.NewServer(logger, router, storageHelper)

	go serv.Run(cfg)

	// Gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("Got signal", zap.String("signal", sign.String()))
	log.Info("Shutting down...")
}

func checkLogFiles(logger *logging.DynamicLogger) {
	zipper := logging.InitZipper(logger.Files)

	for {
		zipper.ZipLogFiles(logger.GetLogger(moduleName), int64(1*1024*1024), "2006-01-02_15:04:05")

		logger.GetLogger(moduleName).Info("Sleeping...")
		time.Sleep(time.Second * 60 * 10)
	}
}

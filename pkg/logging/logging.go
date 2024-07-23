package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"just_for_fun/internal/config"
	"os"
)

var callerKey = "call"

type DynamicLogger struct {
	Files       []*os.File
	modules     map[string]*zap.Logger
	errorCore   zapcore.Core
	defaultCore zapcore.Core
	logsDir     string
	rootLogger  *zap.Logger
}

func InitLogger(cfg *config.Config) *DynamicLogger {
	logsDir := cfg.RootPath + "/" + cfg.Log.Dir
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		errMKDir := os.Mkdir(logsDir, 0755)
		if errMKDir != nil {
			panic(errMKDir)
		}
	}

	newLogger := &DynamicLogger{
		modules: make(map[string]*zap.Logger),
		Files:   []*os.File{},
		logsDir: logsDir,
	}

	newLogger.setErrorCore()
	newLogger.setDefaultCore()

	return newLogger
}

func (l *DynamicLogger) AddModule(moduleName string) {
	if _, ok := l.modules[moduleName]; ok {
		return
	}

	path := l.logsDir + "/" + moduleName + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.modules[moduleName].Error(err.Error(), zap.String("module", moduleName))
	}
	l.Files = append(l.Files, file)

	fileLog := zapcore.AddSync(file)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.CallerKey = callerKey
	moduleEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var moduleCore zapcore.Core

	moduleCore = zapcore.NewCore(moduleEncoder, fileLog, l.lowPriority())

	cores := []zapcore.Core{moduleCore, l.errorCore, l.defaultCore}

	moduleLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	l.modules[moduleName] = moduleLogger
}

func (l *DynamicLogger) GetLogger(moduleName string) *zap.Logger {
	if logger, ok := l.modules[moduleName]; ok {
		logger = logger.Named(moduleName)
		return logger
	}
	return nil
}

func (l *DynamicLogger) Sync() {
	for _, module := range l.modules {
		module.Sync()
	}
	for _, file := range l.Files {
		file.Close()
	}
}

func (l *DynamicLogger) setErrorCore() {
	path := l.logsDir + "/" + zapcore.ErrorLevel.String() + ".log"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	l.Files = append(l.Files, file)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.CallerKey = callerKey
	moduleEncoder := zapcore.NewJSONEncoder(encoderConfig)

	errorCore := zapcore.NewCore(moduleEncoder, zapcore.AddSync(file), l.highPriority())
	l.errorCore = errorCore
}

func (l *DynamicLogger) setDefaultCore() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.CallerKey = callerKey
	moduleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	defaultCore := zapcore.NewCore(moduleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	l.defaultCore = defaultCore
}

func (l *DynamicLogger) highPriority() zap.LevelEnablerFunc {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	return highPriority
}

func (l *DynamicLogger) lowPriority() zap.LevelEnablerFunc {
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	return lowPriority
}

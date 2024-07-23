package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	DB       DBConfig     `yaml:"db"`
	Log      LoggerConfig `yaml:"log"`
	RootPath string
	Swagger  SwaggerConfig `yaml:"swagger"`
	Server   ServerConfig  `yaml:"server"`
}

type ServerConfig struct {
	Host    string `yaml:"host" env-required:"true"`
	Port    string `yaml:"port" env-required:"true"`
	BaseURL string `yaml:"base_url" env-required:"true"`
}

type SwaggerConfig struct {
	URL      string `yaml:"url" env-required:"true"`
	Endpoint string `yaml:"endpoint" env-required:"true"`
}

type LoggerConfig struct {
	Dir string `yaml:"dir" env-required:"true"`
}

type DBConfig struct {
	MigrationDirName string       `yaml:"migration_dir_name" env-required:"true"`
	Driver           string       `yaml:"driver" env-required:"true"`
	Password         string       `yaml:"password" json:"-"`
	SslMode          string       `yaml:"ssl_mode" env-required:"true"`
	Pool             DBPoolConfig `yaml:"pool"`
	Host             string
	Port             string
	User             string
	DbName           string
}

type DBPoolConfig struct {
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

func (d *DBConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		d.Host, d.Port, d.User, d.DbName, d.Password, d.SslMode,
	)
}

func Load(levelsUp int) *Config {
	mustLoadEnvConfig()

	var cfgEnv EnvConfig

	err := cleanenv.ReadEnv(&cfgEnv)
	if err != nil {
		panic(err)
	}

	rootPath := getRootPath(levelsUp)

	pathToCfg := getPath(rootPath, cfgEnv.Dir, cfgEnv.FileName)

	return mustLoadCfg(rootPath, pathToCfg, &cfgEnv.DB)
}

func mustLoadCfg(
	rootPath string,
	pathToCfg string,
	db *EnvDBConfig,
) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(pathToCfg, &cfg)
	if err != nil {
		panic(err)
	}

	cfg.RootPath = rootPath

	addEnvInConfig(&cfg, db)

	return &cfg
}

func createPath(path string, fileName string) {
	_, err := os.Stat(path)
	dir := path
	if fileName != "" {
		dir = filepath.Dir(path)
	}
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getPath(rootPath string, dir string, fileName string) string {
	path := filepath.Join(rootPath, dir, fileName)
	createPath(path, fileName)
	return path
}

func getRootPath(levelsUp int) string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Failed to get root path")
	}

	parentPath := filename
	for i := 0; i < levelsUp; i++ {
		parentPath = filepath.Dir(parentPath)
	}
	return parentPath
}

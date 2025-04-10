package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	GRPCServer  `yaml:"grpc_server"`
	MinIOClient
	Database `yaml:"database"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"10m"`

	AppSecret string
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type GRPCServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:50051"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	User    string `yaml:"user"`
	Name    string `yaml:"name"`
	SSLMode string `yaml:"ssl_mode"`

	Password string
}

type MinIOClient struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func MustLoadEnv() *Config {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found", err)
	}

	// Получаем путь к YAML-конфигу из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	// Читаем YAML-конфиг
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// Получаем обязательные переменные окружения
	cfg.AppSecret = os.Getenv("APP_SECRET")
	if cfg.AppSecret == "" {
		log.Fatal("empty APP_SECRET")
	}

	cfg.Database.User = os.Getenv("DB_USER")
	if cfg.Database.User == "" {
		log.Fatal("empty DB_USER")
	}

	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	if cfg.Database.Password == "" {
		log.Fatal("empty DB_PASSWORD")
	}

	// MinIO
	cfg.MinIOClient.Endpoint = os.Getenv("MINIO_ENDPOINT")
	if cfg.MinIOClient.Endpoint == "" {
		log.Fatal("empty MINIO_ENDPOINT")
	}

	cfg.MinIOClient.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	if cfg.MinIOClient.AccessKey == "" {
		log.Fatal("empty MINIO_ACCESS_KEY")
	}

	cfg.MinIOClient.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	if cfg.MinIOClient.SecretKey == "" {
		log.Fatal("empty MINIO_SECRET_KEY")
	}

	cfg.MinIOClient.UseSSL = false

	return &cfg
}

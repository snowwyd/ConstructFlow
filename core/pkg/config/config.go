package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `env:"ENV"`
	GRPCAddress string `env:"GRPC_ADDRESS"`

	SecretKeys SecretKeys
	HTTPServer HTTPServer
	Database   Database
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}

type SecretKeys struct {
	KeyJwt   string        `env:"KeyJwt"`
	KeyGrpc  string        `env:"KeyGrpc"`
	TokenTTL time.Duration `env:"TokenTTL"`
}

type HTTPServer struct {
	Address string `env:"HTTP_ADDERSS"`
}

func MustLoadEnv() *Config {
	envFile := flag.String("env", "dev", "Env to use dev or prod")
	flag.Parse()

	var envPath string
	println(envFile)
	switch *envFile {
	case "dev":
		envPath = ".env_dev"
	case "prod":
		envPath = ".env_prod"
	default:
		envPath = ".env_dev"
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning .env file not found:")
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning .env file not found: %v", err)
		}
	}

	cfg := &Config{
		GRPCAddress: getEnvParams("GRPC_ADDRESS", "GRPC_ADDRESS"),
		SecretKeys: SecretKeys{
			KeyJwt:   getEnvParams("KeyJwt", "0"),
			KeyGrpc:  getEnvParams("KeyGrpc", "0"),
			TokenTTL: getTimeEnvParams("TokenTTL", "12h"),
		},
		HTTPServer: HTTPServer{
			Address: getEnvParams("HTTP_ADDRESS", "0.0.0.0:8080"),
		},
		Database: Database{
			Host:     getEnvParams("DB_HOST", "localhost"),
			Port:     getIntEnvParams("DB_PORT", 5432),
			Name:     getEnvParams("DB_NAME", "constructflow"),
			User:     getEnvParams("DB_USER", "postgres"),
			Password: getEnvParams("DB_PASSWORD", "postgres"),
		},
	}

	return cfg
}

func getEnvParams(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvParams(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Could not parse %s as integer, using default value: %v", key, err)
		return defaultValue
	}

	return value
}

func getTimeEnvParams(key string, defaultValue string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		val, err := time.ParseDuration(defaultValue)
		if err != nil {
			log.Printf("Warning: Could not parse %s as duration, using default value: %v", key, err)
			return 0
		}
		return val
	}

	duration, _ := time.ParseDuration(value)
	return duration
}

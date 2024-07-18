package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Env  string    `yaml:"env" env-default:"local"`
	Port string    `yaml:"port" envDefault:"8080"`
	DB   DBConfig  `yaml:"db"`
	API  APIConfig `yaml:"api"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type APIConfig struct {
	ExternalURL string `yaml:"api_url"`
}

func MustLoad() Config {
	absPath, err := filepath.Abs(".env/")

	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		panic("config file does not exist: " + absPath)
	}

	if err := godotenv.Load(absPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var cfg Config
	cfg.Port = os.Getenv("PORT")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Username = os.Getenv("DB_USERNAME")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.DBName = os.Getenv("DB_NAME")
	cfg.DB.SSLMode = os.Getenv("SSL_MODE")
	cfg.API.ExternalURL = os.Getenv("API_URL")

	return cfg
}

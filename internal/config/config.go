package config

import (
	"log"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type AppConfig struct {
	AppEnv      string `koanf:"APP_ENV"`
	AppPort     string `koanf:"APP_PORT"`
	AppGrpcPort string `koanf:"APP_GRPC_PORT"`
	DbURI       string `koanf:"DB_URI"`
}

func SetupConfig() *AppConfig {
	// Instance koanf
	k := koanf.New(".")

	// Read .env file
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		log.Println("error loading .env file : ", err)
	}

	// Read OS
	if err := k.Load(env.Provider("", ".", func(s string) string {
		return s
	}), nil); err != nil {
		log.Println("error load env provider ", err)
	}

	// Default value
	config := AppConfig{
		AppEnv:      "development",
		AppPort:     "8080",
		AppGrpcPort: "9090",
		DbURI:       "user:password@tcp(127.0.0.1:3306)/ewallet?charset=utf8mb4&parseTime=True&loc=Local",
	}

	// Unmarshal to struct go
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("failed mapping config : %v", err)
	}

	return &config
}

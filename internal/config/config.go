package config

import (
	"fmt"
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
	DbHost      string `koanf:"DB_HOST"`
	DbPort      string `koanf:"DB_PORT"`
	DbUser      string `koanf:"DB_USER"`
	DbPassword  string `koanf:"DB_PASSWORD"`
	DbName      string `koanf:"DB_NAME"`
	UmsGrpcHost string `koanf:"UMS_GRPC_HOST"`
}

func (a *AppConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", a.DbUser, a.DbPassword, a.DbHost, a.DbPort, a.DbName)
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
		AppPort:     "8081",
		AppGrpcPort: "9091",
		DbHost:      "127.0.0.1",
		DbPort:      "3306",
		DbUser:      "root",
		DbPassword:  "password",
		DbName:      "ewallet_wallet",
		UmsGrpcHost: "localhost:9090",
	}

	// Unmarshal to struct go
	if err := k.Unmarshal("", &config); err != nil {
		log.Fatalf("failed mapping config : %v", err)
	}

	return &config
}

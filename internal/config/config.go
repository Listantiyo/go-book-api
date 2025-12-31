package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName   string
	AppPort   string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	JWTSecret string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetDefault("APP_NAME", "App Name")
	viper.SetDefault("APP_PORT", "8080")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASS", "123123")
	viper.SetDefault("DB_NAME", "book_api")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "secret")

	if err := viper.ReadInConfig(); err != nil {
		 log.Println("No .env file found, using environment variables or defaults")
	}else{
		log.Println("✅ Configuration loaded successfully.")
	}


	return &Config{
		AppName: viper.GetString("APP_NAME"),
		AppPort: viper.GetString("APP_PORT"),

		DBHost: viper.GetString("DB_HOST"),
		DBPort: viper.GetString("DB_PORT"),
		DBUser: viper.GetString("DB_USER"),
		DBPass: viper.GetString("DB_PASS"),
		DBName: viper.GetString("DB_NAME"),
		DBSSLMode: viper.GetString("DB_SSLMODE"),

		JWTSecret: viper.GetString("JWT_SECRET"),
	}
}
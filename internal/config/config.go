package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslmode  string

	REDISHost     string
	REDISPort     string
	REDISPassword string
	REDISDb       string

	ServerHost string
	ServerPort string
	ServerEnv  string
}

var DbConfig *Config

func ConfigDb() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	setDefaults()

	DbConfig = &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSslmode:  viper.GetString("DB_SSLMODE"),

		REDISHost: viper.GetString("REDIS_HOST"),
		REDISPort: viper.GetString("REDIS_PORT"),
		REDISPassword: viper.GetString("REDIS_PASSWORd"),
		REDISDb: viper.GetString("REDIS_DB"),

		ServerHost: viper.GetString("SERVER_HOST"),
		ServerPort: viper.GetString("SERVER_PORT"),
		ServerEnv:  viper.GetString("SERVER_ENV"),
	}

	log.Println("Configuration loaded successfully")
}

func setDefaults() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "course_online_db")
	viper.SetDefault("DB_SSLMODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", "0")

	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "7070")
	viper.SetDefault("SERVER_ENV", "development")
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

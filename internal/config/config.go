package config

import (
	"log"

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

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFromEmail    string
	SMTPFromName     string

	OTPExpiryMinutes int
	OTPLength        int
}

var DbConfig *Config

func ConfigDb() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	setDefaults()

	DbConfig = &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBSslmode:  viper.GetString("DB_SSLMODE"),

		REDISHost:     viper.GetString("REDIS_HOST"),
		REDISPort:     viper.GetString("REDIS_PORT"),
		REDISPassword: viper.GetString("REDIS_PASSWORD"),
		REDISDb:       viper.GetString("REDIS_DB"),

		ServerHost: viper.GetString("SERVER_HOST"),
		ServerPort: viper.GetString("SERVER_PORT"),
		ServerEnv:  viper.GetString("SERVER_ENV"),

		SMTPHost:     viper.GetString("SMTP_HOST"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPassword: viper.GetString("SMTP_PASSWORD"),
		SMTPFromEmail:    viper.GetString("SMTP_FROM_EMAIL"),
		SMTPFromName:     viper.GetString("SMTP_FROM_NAME"),

		OTPExpiryMinutes: viper.GetInt("OTP_EXPIRY_MINUTES"),
		OTPLength:        viper.GetInt("OTP_LENGTH"),
	}

	log.Printf("[CONFIG] SMTP Host: %s", DbConfig.SMTPHost)
	log.Printf("[CONFIG] SMTP Port: %d", DbConfig.SMTPPort)
	log.Printf("[CONFIG] SMTP User: %s", DbConfig.SMTPUser)
	log.Printf("[CONFIG] SMTP Password length: %d chars", len(DbConfig.SMTPPassword))
	log.Printf("[CONFIG] SMTP From Email: %s", DbConfig.SMTPFromEmail)
	
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

	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USER", "")
	viper.SetDefault("SMTP_PASSWORD", "")
	viper.SetDefault("SMTP_FROM_EMAIL", "")
	viper.SetDefault("SMTP_FROM_NAME", "Course Online")

	viper.SetDefault("OTP_EXPIRY_MINUTES", 5)
	viper.SetDefault("OTP_LENGTH", 6)
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
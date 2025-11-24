package config

import (
	"log"

	"github.com/spf13/viper"
)

type SMTPConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	FromEmail   string
	FromName    string
}

var SmtpConfig *SMTPConfig

func LoadSMTPConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	setSMTPDefaults()

	SmtpConfig = &SMTPConfig{
		Host:      viper.GetString("SMTP_HOST"),
		Port:      viper.GetInt("SMTP_PORT"),
		User:      viper.GetString("SMTP_USER"),
		Password:  viper.GetString("SMTP_PASSWORD"),
		FromEmail: viper.GetString("SMTP_FROM_EMAIL"),
		FromName:  viper.GetString("SMTP_FROM_NAME"),
	}

	log.Printf("[SMTP CONFIG] Host: %s", SmtpConfig.Host)
	log.Printf("[SMTP CONFIG] Port: %d", SmtpConfig.Port)
	log.Printf("[SMTP CONFIG] User: %s", SmtpConfig.User)
	log.Printf("[SMTP CONFIG] Pass length: %d chars", len(SmtpConfig.Password))
	log.Printf("[SMTP CONFIG] From Email: %s", SmtpConfig.FromEmail)
}

func setSMTPDefaults() {
	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USER", "")
	viper.SetDefault("SMTP_PASSWORD", "")
	viper.SetDefault("SMTP_FROM_EMAIL", "")
	viper.SetDefault("SMTP_FROM_NAME", "Course Online")
}

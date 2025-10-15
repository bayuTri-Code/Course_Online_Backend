package config

import (
	"log"

	"github.com/spf13/viper"
)

type MinioConfiguration struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

var MinioConfig *MinioConfiguration

func InitMinioConfig() {
	MinioConfig = &MinioConfiguration{
		Endpoint:  viper.GetString("MINIO_ENDPOINT"),
		AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
		SecretKey: viper.GetString("MINIO_SECRET_KEY"),
		Bucket:    viper.GetString("MINIO_BUCKET"),
		UseSSL:    viper.GetBool("MINIO_USE_SSL"),
	}

	log.Println("[Config] MinIO configuration loaded successfully")
	log.Printf("[Config] MinIO Endpoint: %s", MinioConfig.Endpoint)
	log.Printf("[Config] MinIO Bucket: %s", MinioConfig.Bucket)
	log.Printf("[Config] MinIO UseSSL: %v", MinioConfig.UseSSL)
}


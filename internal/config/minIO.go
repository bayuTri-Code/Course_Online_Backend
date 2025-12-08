package config

import (
	"log"

	"github.com/spf13/viper"
)

type MinioConfiguration struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Bucket        string
	UseSSL        bool
	ImageENDPOINT string
	MinSize       int64
	MaxSize       int64
}

var MinioConfig *MinioConfiguration

func InitMinioConfig() {
	MinioConfig = &MinioConfiguration{
		Endpoint:      viper.GetString("MINIO_ENDPOINT"),
		AccessKey:     viper.GetString("MINIO_ACCESS_KEY"),
		SecretKey:     viper.GetString("MINIO_SECRET_KEY"),
		Bucket:        viper.GetString("MINIO_BUCKET"),
		UseSSL:        viper.GetBool("MINIO_USE_SSL"),
		ImageENDPOINT: viper.GetString("IMAGE_ENDPOINT"),
		MinSize:       viper.GetInt64("MINIO_MIN_UPLOAD_SIZE"),
		MaxSize:       viper.GetInt64("MINIO_MAX_UPLOAD_SIZE"),
	}

	log.Println("[Config] MinIO configuration loaded successfully")
	log.Printf("[Config] MinIO Endpoint: %s", MinioConfig.Endpoint)
	log.Printf("[Config] MinIO Bucket: %s", MinioConfig.Bucket)
	log.Printf("[Config] MinIO UseSSL: %v", MinioConfig.UseSSL)
	log.Printf("[Config] MinIO MinSize: %d bytes", MinioConfig.MinSize)
	log.Printf("[Config] MinIO MaxSize: %d bytes", MinioConfig.MaxSize)
}

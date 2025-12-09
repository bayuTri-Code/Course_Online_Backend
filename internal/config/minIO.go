package config

import (
	"log"
	"strconv"
	"strings"

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
	minSizeStr := viper.GetString("MINIO_MIN_UPLOAD_SIZE")
	minSize, err := parseSize(minSizeStr)
	if err != nil {
		log.Fatalf("[Config] Invalid MINIO_MIN_UPLOAD_SIZE: %s", minSizeStr)
	}

	maxSizeStr := viper.GetString("MINIO_MAX_UPLOAD_SIZE")
	maxSize, err := parseSize(maxSizeStr)
	if err != nil {
		log.Fatalf("[Config] Invalid MINIO_MAX_UPLOAD_SIZE: %s", maxSizeStr)
	}

	MinioConfig = &MinioConfiguration{
		Endpoint:      viper.GetString("MINIO_ENDPOINT"),
		AccessKey:     viper.GetString("MINIO_ACCESS_KEY"),
		SecretKey:     viper.GetString("MINIO_SECRET_KEY"),
		Bucket:        viper.GetString("MINIO_BUCKET"),
		UseSSL:        viper.GetBool("MINIO_USE_SSL"),
		ImageENDPOINT: viper.GetString("IMAGE_ENDPOINT"),
		MinSize:       minSize, 
		MaxSize:       maxSize,
	}

	log.Println("[Config] MinIO configuration loaded successfully")
	log.Printf("[Config] MinIO Endpoint: %s", MinioConfig.Endpoint)
	log.Printf("[Config] MinIO Bucket: %s", MinioConfig.Bucket)
	log.Printf("[Config] MinIO UseSSL: %v", MinioConfig.UseSSL)
	log.Printf("[Config] MinIO MinSize: %d bytes", MinioConfig.MinSize)
	log.Printf("[Config] MinIO MaxSize: %d bytes", MinioConfig.MaxSize)
}

func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	if strings.HasSuffix(sizeStr, "KB") {
		value, err := strconv.Atoi(strings.TrimSuffix(sizeStr, "KB"))
		if err != nil {
			return 0, err
		}
		return int64(value) * 1024, nil
	}

	if strings.HasSuffix(sizeStr, "MB") {
		value, err := strconv.Atoi(strings.TrimSuffix(sizeStr, "MB"))
		if err != nil {
			return 0, err
		}
		return int64(value) * 1024 * 1024, nil
	}

	if strings.HasSuffix(sizeStr, "B") {
		value, err := strconv.Atoi(strings.TrimSuffix(sizeStr, "B"))
		if err != nil {
			return 0, err
		}
		return int64(value), nil
	}

	value, err := strconv.Atoi(sizeStr)
	return int64(value), err
}

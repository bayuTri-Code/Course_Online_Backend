package database

import (
	"context"
	
	"log"

	"course_online_backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func MinioConn() *minio.Client {
	cfg := config.MinioConfig

	if cfg == nil {
		log.Fatal("[MinIO] Configuration is nil. Please call config.InitMinioConfig() first")
		return nil
	}

	if cfg.Endpoint == "" {
		log.Fatal("[MinIO] Endpoint is empty")
		return nil
	}

	log.Printf("[MinIO] Attempting to connect to endpoint: %s", cfg.Endpoint)
	log.Printf("[MinIO] Using credentials - AccessKey: %s", cfg.AccessKey)

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Printf("[MinIO] Failed to create client: %v", err)
		return nil
	}

	log.Printf("[MinIO] Client created successfully")

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		log.Printf("[MinIO] Failed to check bucket existence: %v", err)
		log.Printf("[MinIO] This might be a connection issue. Please verify:")
		return nil
	}

	if !exists {
		log.Printf("[MinIO] Bucket '%s' does not exist, creating...", cfg.Bucket)
		
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("[MinIO] Failed to create bucket '%s': %v", cfg.Bucket, err)
			return nil
		}
		
		log.Printf("[MinIO] Bucket '%s' created successfully", cfg.Bucket)
	} else {
		log.Printf("[MinIO] Bucket '%s' already exists", cfg.Bucket)
	}

	log.Println("CONNECTION SUCCESSFUL ")
	
	MinioClient = client
	return client
}

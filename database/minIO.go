// File: internal/database/minio.go
package database

import (
	"context"
	"fmt"
	"log"

	"course_online_backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func MinioConn() *minio.Client {
	cfg := config.MinioConfig

	if cfg == nil {
		log.Fatal(" Configuration is nil. Please call config.InitMinioConfig() first")
		return nil
	}

	if cfg.Endpoint == "" {
		log.Fatal(" Endpoint is empty")
		return nil
	}

	log.Printf(" Attempting to connect to endpoint: %s", cfg.Endpoint)
	log.Printf(" Using credentials - AccessKey: %s", cfg.AccessKey)

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return nil
	}

	log.Printf("Client created successfully")

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		log.Printf("Failed to check bucket existence: %v", err)
		log.Printf("This might be a connection issue. Please verify:")
		return nil
	}

	if !exists {
		log.Printf("Bucket '%s' does not exist, creating...", cfg.Bucket)
		
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf(" Failed to create bucket '%s': %v", cfg.Bucket, err)
			return nil
		}
		
		log.Printf("Bucket '%s' created successfully", cfg.Bucket)
	} else {
		log.Printf("Bucket '%s' already exists", cfg.Bucket)
	}

	if err := setPublicBucketPolicy(client, cfg.Bucket); err != nil {
		log.Printf("Warning: Failed to set public policy: %v", err)
		log.Printf("Files may not be publicly accessible")
	} else {
		log.Printf(" Bucket policy set to public read successfully")
	}

	log.Println("CONNECTION SUCCESSFUL")
	
	MinioClient = client
	return client
}

func setPublicBucketPolicy(client *minio.Client, bucketName string) error {
	ctx := context.Background()

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)

	err := client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}
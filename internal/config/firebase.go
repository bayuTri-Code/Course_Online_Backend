package config

import (
	"context"
	"log"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)


func InitFirebase() *firebase.App {
	opt := option.WithCredentialsFile("internal/config/serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing firebase: %v", err)
	}
	log.Println("Firebase initialized successfully")
	return app
}

package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
)

type contextKey string

func FirebaseAuth(app *firebase.App) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			idToken := strings.Replace(authHeader, "Bearer ", "", 1)
			client, err := app.Auth(context.Background())
			if err != nil {
				http.Error(w, "Error initializing auth", http.StatusInternalServerError)
				return
			}

			token, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextKey("userUID"), token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

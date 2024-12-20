package middleware

import (
	"context"
	"net/http"
	"strings"

	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

type contextKey string

const AdminContextKey contextKey = "admin"

func AuthMiddleware(jwtUtils *utils.JWTUtils, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization Header Missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		claims, err := jwtUtils.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AdminContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

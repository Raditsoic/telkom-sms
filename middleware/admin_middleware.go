package middleware

import (
	"context"
	"net/http"
	"strings"

	"gtihub.com/raditsoic/telkom-storage-ms/utils"
)

const AdminContextKey = "admin"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		claims, err := utils.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AdminContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

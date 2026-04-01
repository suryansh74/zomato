package middleware

import (
	"net/http"
)

// The shared secret (In production, load this from .env!)
const InternalSecretKey = "super_secret_zomato_internal_key_2026"

func InternalAPIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Look for the VIP pass in the headers
		clientKey := r.Header.Get("X-Internal-Key")

		if clientKey != InternalSecretKey {
			http.Error(w, "Unauthorized Internal Service Call", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

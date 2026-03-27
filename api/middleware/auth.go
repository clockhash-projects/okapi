package middleware

import (
	"crypto/subtle"
	"net/http"
	"okapi/internal/config"
)

func Auth(cfg *config.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				// Also check Authorization: Bearer <token>
				authHeader := r.Header.Get("Authorization")
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					apiKey = authHeader[7:]
				}
			}

			if apiKey == "" {
				http.Error(w, "Unauthorized: API Key missing", http.StatusUnauthorized)
				return
			}

			valid := false
			for _, key := range cfg.APIKeys {
				if subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) == 1 {
					valid = true
					break
				}
			}

			if !valid {
				http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package auth0

import (
	"net/http"

	// "github.com/golang-jwt/jwt/v5"
)

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: extract and validate token
		
	})
}
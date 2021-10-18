package ssojwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func MakeAccessTokenMiddleware(config SSOConfig, key string) func(nextHandler http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			AuthorizationMap := strings.Split(authorization, " ")
			if len(AuthorizationMap) < 2 {
				nextHandler.ServeHTTP(w, r)
				return
			}
			tokenString := AuthorizationMap[1]
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.AccessTokenSecretKey), nil
			})
			if err != nil {
				nextHandler.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()

			claims, ok := token.Claims.(jwt.MapClaims)
			if ok && token.Valid {
				ctx = context.WithValue(ctx, key, claims)
			}
			nextHandler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
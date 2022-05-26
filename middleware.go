package ssojwt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func MakeAccessTokenMiddleware(config SSOConfig, key string) func(nextHandler http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			AuthorizationMap := strings.Split(authorization, " ")
			if len(AuthorizationMap) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "{\"error\": \"invalid_token\"}")
				return
			}
			tokenString := AuthorizationMap[1]
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return []byte(config.AccessTokenSecretKey), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "{\"error\": \"invalid_token\"}")
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

func MakeRefreshTokenMiddleware(config SSOConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		AuthorizationMap := strings.Split(authorization, " ")
		if len(AuthorizationMap) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "{\"error\": \"invalid_token\"}")
			return
		}
		tokenString := AuthorizationMap[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.RefreshTokenSecretKey), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "{\"error\": \"invalid_token\"}")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			jurusan := claims["jurusan"].(map[string]interface{})

			newClaims := ServiceResponse{
				AuthenticationSuccess: AuthenticationSuccess{
					User: claims["user"].(string),
					Attributes: Attributes{
						Nama: claims["nama"].(string),
						Npm:  claims["npm"].(string),
						Jurusan: Jurusan{
							Faculty:      jurusan["faculty"].(string),
							ShortFaculty: jurusan["shortFaculty"].(string),
							Major:        jurusan["major"].(string),
							Program:      jurusan["program"].(string),
						},
					},
				},
			}
			accessToken, err := CreateAccessToken(config, newClaims)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "{\"error\": \"internal_server_error\"}")
				return
			}

			refreshToken, err := CreateRefreshToken(config, newClaims)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "{\"error\": \"internal_server_error\"}")
				return
			}
			res := LoginResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				Fakultas:     nil,
			}

			w.WriteHeader(http.StatusOK)
			resJson, _ := json.Marshal(res)
			fmt.Fprintf(w, "%s", resJson)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{\"error\": \"invalid_token\"}")
		return

	})

}

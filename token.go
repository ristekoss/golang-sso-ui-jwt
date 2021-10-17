package ssojwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SSOJwtClaim struct {
	Nama    string  `json:"nama"`
	User    string  `json:"user"`
	Npm     string  `json:"npm"`
	Jusuran Jurusan `json:"jurusan"`
	jwt.StandardClaims
}

func CreateAccessToken(config SSOConfig, ssoResponse ServiceResponse) (token string, err error) {
	claims := &SSOJwtClaim{
		Nama:    ssoResponse.AuthenticationSuccess.Attributes.Nama,
		User:    ssoResponse.AuthenticationSuccess.User,
		Npm:     ssoResponse.AuthenticationSuccess.Attributes.Npm,
		Jusuran: ssoResponse.AuthenticationSuccess.Attributes.Jusuran,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.AccessTokenExpireTime).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = rawToken.SignedString([]byte(config.AccessTokenSecretKey))
	if err != nil {
		err = fmt.Errorf("token signing error: %w", err)
		return
	}
	return
}

func CreateRefreshToken(config SSOConfig, ssoResponse ServiceResponse) (token string, err error) {
	claims := &SSOJwtClaim{
		Nama:    ssoResponse.AuthenticationSuccess.Attributes.Nama,
		User:    ssoResponse.AuthenticationSuccess.User,
		Npm:     ssoResponse.AuthenticationSuccess.Attributes.Npm,
		Jusuran: ssoResponse.AuthenticationSuccess.Attributes.Jusuran,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.RefreshTokenExpireTime).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = rawToken.SignedString([]byte(config.RefreshTokenSecretKey))
	if err != nil {
		return "", fmt.Errorf("token signing error: %w", err)
	}
	return
}

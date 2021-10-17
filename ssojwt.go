package ssojwt

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type SSOConfig struct {
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
	AccessTokenSecretKey   string
	RefreshTokenSecretKey  string
	ServiceUrl             string
	OriginUrl              string
	CasURL                 string
}

func MakeSSOConfig(accessTokenExpireTime, refreshTokenExpireTime time.Duration, accessTokenSecretKey, refreshTokenSecretKey, serviceUrl, originUrl string) SSOConfig {
	return SSOConfig{
		AccessTokenExpireTime:  accessTokenExpireTime,
		RefreshTokenExpireTime: refreshTokenExpireTime,
		AccessTokenSecretKey:   accessTokenSecretKey,
		RefreshTokenSecretKey:  refreshTokenSecretKey,
		ServiceUrl:             serviceUrl,
		OriginUrl:              originUrl,
		CasURL:                 "https://sso.ui.ac.id/cas2/",
	}
}

func LoginCreator(config SSOConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ticket := r.URL.Query().Get("ticket")
		bodyBytes, err := ValidatTicket(config, ticket)
		if err != nil {
			log.Printf("error when cheking ticket: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		model, err := Unmarshal(bodyBytes)
		if err != nil {
			log.Printf("error in unmarshaling: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Printf("error in reading json: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := MakeLoginResponse(config, model)
		if err != nil {
			log.Printf("error in creating token: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tempData := DataRender{
			LoginResponse: res,
			OriginUrl:     config.OriginUrl,
		}

		abs, _ := filepath.Abs("../static/wait.html")
		tmpl, err := template.ParseFiles(abs)
		err = tmpl.Execute(w, tempData)
		if err != nil {
			log.Printf("error in execute: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

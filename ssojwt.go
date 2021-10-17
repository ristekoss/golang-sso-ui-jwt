package ssojwt

import (
	"io/ioutil"
	"log"
	"net/http"
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

func LoginCreator(config SSOConfig, errorLogger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	if errorLogger == nil {
		errorLogger = log.New(ioutil.Discard, "Error: ", log.Ldate|log.Ltime)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		ticket := r.URL.Query().Get("ticket")
		bodyBytes, err := ValidatTicket(config, ticket)
		if err != nil {
			errorLogger.Printf("error when cheking ticket: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		model, err := Unmarshal(bodyBytes)
		if err != nil {
			errorLogger.Printf("error in unmarshaling: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := MakeLoginResponse(config, model)
		if err != nil {
			errorLogger.Printf("error in creating token: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, dataRender, err := MakeTemplate(config, res)
		if err != nil {
			errorLogger.Printf("error in parsing template: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = tmpl.Execute(w, dataRender)
		if err != nil {
			errorLogger.Printf("error in render template: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

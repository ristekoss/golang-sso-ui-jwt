package ssojwt

import (
	"fmt"
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
		res, err := LoginRequestHandler(ticket, config)
		if err != nil {
			errorLogger.Printf("error in pasing sso request: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = TemplateRenderHandler(res, config, w)
		if err != nil {
			errorLogger.Printf("error in render template: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func LoginRequestHandler(ticket string, config SSOConfig) (res LoginResponse, err error) {
	bodyBytes, err := ValidatTicket(config, ticket)
	if err != nil {
		err = fmt.Errorf("error when cheking ticket: %w", err)
		return
	}

	model, err := Unmarshal(bodyBytes)
	if err != nil {
		err = fmt.Errorf("error in unmarshaling: %w", err)
		return
	}

	res, err = MakeLoginResponse(config, model)
	if err != nil {
		err = fmt.Errorf("error in creating token: %w", err)
	}
	return
}

func TemplateRenderHandler(data interface{}, config SSOConfig, w http.ResponseWriter) (err error) {
	tmpl, dataRender, err := MakeTemplate(config, data)
	if err != nil {
		err = fmt.Errorf("error in making template: %w", err)
		return
	}

	err = tmpl.Execute(w, dataRender)
	if err != nil {
		err = fmt.Errorf("error in parsing template: %w", err)
	}
	return
}

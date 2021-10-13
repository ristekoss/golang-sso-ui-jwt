package ssojwt

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SSOConfig struct {
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
	AccessTokenSecretKey   string
	RefreshTokenSecretKey  string
	ServiceUrl             string
	OriginUrl              string
}

func CreateSSOConfig(accessTokenExpireTime, refreshTokenExpireTime time.Duration, accessTokenSecretKey, refreshTokenSecretKey, serviceUrl, originUrl string) SSOConfig {
	return SSOConfig{
		AccessTokenExpireTime:  accessTokenExpireTime,
		RefreshTokenExpireTime: refreshTokenExpireTime,
		AccessTokenSecretKey:   accessTokenSecretKey,
		RefreshTokenSecretKey:  refreshTokenSecretKey,
		ServiceUrl:             serviceUrl,
		OriginUrl:              originUrl,
	}
}

func LoginCreator(config SSOConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		casURL := "https://sso.ui.ac.id/cas2/"
		ticket := r.URL.Query().Get("ticket")

		url := fmt.Sprintf("%sserviceValidate?ticket=%s&service=%s", casURL, ticket, config.ServiceUrl)
		resp, err := http.Get(url)
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error when cheking ticket: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var model ServiceResponse
		xml.Unmarshal(bodyBytes, &model)
		data, err := readOrgcode()
		if err != nil {
			log.Printf("error in reading json: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		model.AuthenticationSuccess.Attributes.Jusuran = data[model.AuthenticationSuccess.Attributes.Kd_org]
		accessToken, _ := createToken(config, model)
		refreshToken, _ := createRefreshToken(config, model)
		res := LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Nama:         model.AuthenticationSuccess.Attributes.Nama,
			Npm:          model.AuthenticationSuccess.Attributes.Npm,
			Fakultas:     model.AuthenticationSuccess.Attributes.Jusuran.Faculty,
		}

		tempData := DataRender{
			LoginResponse: res,
			OriginUrl: config.OriginUrl,
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

func readOrgcode() (data map[string]Jurusan, err error) {
	abs, _ := filepath.Abs("../static/orgcode.json")
	file, _ := ioutil.ReadFile(abs)
	err = json.Unmarshal([]byte(file), &data)
	return
}

func createToken(config SSOConfig, ssoResponse ServiceResponse) (token string, err error) {
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
		return "", fmt.Errorf("token signing error: %w", err)
	}
	return
}

func createRefreshToken(config SSOConfig, ssoResponse ServiceResponse) (token string, err error) {
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

type DataRender struct {
	LoginResponse LoginResponse
	OriginUrl     string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Nama         string `json:"nama,omitempty"`
	Npm          string `json:"npm,omitempty"`
	Fakultas     string `json:"fakultas,omitempty"`
}

type SSOJwtClaim struct {
	Nama    string  `json:"nama"`
	User    string  `json:"user"`
	Npm     string  `json:"npm"`
	Jusuran Jurusan `json:"jurusan"`
	jwt.StandardClaims
}

type ServiceResponse struct {
	XMLName               xml.Name              `xml:"serviceResponse" json:"-"`
	AuthenticationSuccess AuthenticationSuccess `xml:"authenticationSuccess"`
}

type AuthenticationSuccess struct {
	XMLName    xml.Name   `xml:"authenticationSuccess" json:"-"`
	User       string     `xml:"user" json:"user"`
	Attributes Attributes `xml:"attributes" json:"attributes"`
}

type Attributes struct {
	XMLName    xml.Name `xml:"attributes" json:"-"`
	Ldap_cn    string   `xml:"ldap_cn" xml:"ldap_cn"`
	Kd_org     string   `xml:"kd_org" json:"kd_org"`
	Peran_user string   `xml:"peran_user" json:"peran_user"`
	Nama       string   `xml:"nama" json:"nama"`
	Npm        string   `xml:"npm" json:"npm"`
	Jusuran    Jurusan  `json:"jurusan"`
}

type Jurusan struct {
	Faculty      string `json:"faculty"`
	ShortFaculty string `json:"shortFaculty"`
	Major        string `json:"major"`
	Program      string `json:"program"`
}

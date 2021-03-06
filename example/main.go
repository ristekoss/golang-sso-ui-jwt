package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v4"
	ssojwt "github.com/ristekoss/golang-sso-ui-jwt"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		abs, _ := filepath.Abs("./login.html")
		tmpl, err := template.ParseFiles(abs)
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("error in parse: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	config := ssojwt.MakeSSOConfig(time.Hour*168, time.Hour*720, "super secret access", "super secret refresh", "http://localhost:8080/login", "http://localhost:8080/")
	http.HandleFunc("/login", ssojwt.LoginCreator(config, nil))

	middle := ssojwt.MakeAccessTokenMiddleware(config, "user")
	check := middle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value("user")
		if data == nil {
			data = jwt.MapClaims{"npm": "none"}
		}
		user, _ := json.Marshal(data)
		fmt.Fprintf(w, "%s", user)
	}))
	http.Handle("/check", check)

	refresh := ssojwt.MakeRefreshTokenMiddleware(config)
	http.Handle("/refresh", refresh)

	fmt.Println("server started at localhost:8080")
	http.ListenAndServe(":8080", nil)
}

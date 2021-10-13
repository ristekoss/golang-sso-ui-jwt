package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	ssojwt "github.com/RistekCSUI/golang-sso-ui-jwt"
)

func main() {
	config := ssojwt.CreateSSOConfig(time.Hour*168, time.Hour*720, "super secret", "huha huha", "http://localhost:8080/login", "http://localhost:8080/")
	http.HandleFunc("/login", ssojwt.LoginCreator(config))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		abs, _ := filepath.Abs("./test.html")
		tmpl, err := template.ParseFiles(abs)
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("error in parse: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println("server started at localhost:8080")
	http.ListenAndServe(":8080", nil)
}

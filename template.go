package ssojwt

import (
	"html/template"
	"path/filepath"
)

func MakeTemplate(config SSOConfig, res interface{}) (tmpl *template.Template, dataRender DataRender, err error) {
	dataRender = DataRender{
		LoginResponse: res,
		OriginUrl:     config.OriginUrl,
	}

	abs, err := filepath.Abs("../static/wait.html")
	if err != nil {
		return
	}
	tmpl, err = template.ParseFiles(abs)
	return
}

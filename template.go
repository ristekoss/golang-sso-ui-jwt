package ssojwt

import (
	"html/template"
)

func MakeTemplate(config SSOConfig, res interface{}) (tmpl *template.Template, dataRender DataRender, err error) {
	dataRender = DataRender{
		LoginResponse: res,
		OriginUrl:     config.OriginUrl,
	}

	tmpl, err = template.New("wait").Parse(wait)
	return
}

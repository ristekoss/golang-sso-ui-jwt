package ssojwt

import (
	"encoding/json"
	"fmt"
	"html/template"
)

func MakeTemplate(config SSOConfig, res interface{}) (tmpl *template.Template, dataRender DataRender, err error) {
	dataByte, err := json.Marshal(res)
	if err != nil {
		return
	}
	dataString := fmt.Sprintf("%s", dataByte)
	dataRender = DataRender{
		LoginResponse: dataString,
		OriginUrl:     config.OriginUrl,
	}

	tmpl, err = template.New("wait").Parse(wait)
	return
}

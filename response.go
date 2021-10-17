package ssojwt

type DataRender struct {
	LoginResponse interface{}
	OriginUrl     string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Nama         string `json:"nama,omitempty"`
	Npm          string `json:"npm,omitempty"`
	Fakultas     string `json:"fakultas,omitempty"`
}

func MakeLoginResponse(config SSOConfig, model ServiceResponse) (res LoginResponse, err error) {
	accessToken, err := CreateAccessToken(config, model)
	if err != nil {
		return
	}
	refreshToken, err := CreateRefreshToken(config, model)
	if err != nil {
		return
	}

	res = LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Nama:         model.AuthenticationSuccess.Attributes.Nama,
		Npm:          model.AuthenticationSuccess.Attributes.Npm,
		Fakultas:     model.AuthenticationSuccess.Attributes.Jusuran.Faculty,
	}
	return
}

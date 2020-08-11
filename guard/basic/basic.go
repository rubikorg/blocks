package basic

import (
	"encoding/base64"
	"net/http"
	"strings"

	r "github.com/rubikorg/rubik"
)

type BasicGuard struct{}

type config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ResponseType  string `json:"response_type"`
	NoHeaderError string `json:"no_header_error"`
	AuthError     string `json:"auth_error"`
}

const (
	authorizationKey = "Authorization"
	BlockName        = "BasicGuard"
)

var guardConfig config

func (bg BasicGuard) OnAttach(app *r.App) error {
	err := app.Decode("basic_auth", &guardConfig)
	if err != nil {
		return err
	}
	return nil
}

func AuthGuard(req *r.Request) {
	responseType := guardConfig.ResponseType
	if responseType == "" {
		responseType = "text"
	}

	header := req.Raw.Header
	if guardConfig.Username == "" || guardConfig.Password == "" {
		msg := "BasicGuard: requires you to specify username & password inside" +
			" [basic_auth] object of your config"
		req.Throw(500, r.E(msg), r.StringByteTypeMap[responseType])
		return
	}

	aHeader := header.Get(authorizationKey)
	if !strings.Contains(aHeader, "Basic") {
		msg := guardConfig.NoHeaderError
		if msg == "" {
			msg = "This request doesn't have basic authorization key"
		}
		req.Throw(http.StatusUnauthorized, r.E(msg), r.StringByteTypeMap[responseType])
		return
	}

	aHeader = strings.Replace(aHeader, "Basic ", "", 1)
	decoded, err := base64.StdEncoding.DecodeString(aHeader)
	if err != nil {
		req.Throw(http.StatusUnauthorized, err, r.StringByteTypeMap[responseType])
		return
	}

	decData := strings.Split(string(decoded), ":")
	if len(decData) == 0 || len(decData) < 2 {
		req.Throw(http.StatusUnauthorized,
			r.E("Malformed Basic header."), r.StringByteTypeMap[responseType])
		return
	}

	if decData[0] != guardConfig.Username || guardConfig.Password != decData[1] {
		req.Throw(http.StatusUnauthorized,
			r.E("Authorization failure!"), r.StringByteTypeMap[responseType])
		return
	}
}

func init() {
	r.Attach(BlockName, BasicGuard{})
}

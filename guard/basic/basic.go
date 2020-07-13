package guard

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	r "github.com/rubikorg/rubik"
)

type BasicGuard struct {
	config guardConfig
}

type guardConfig struct {
	username string
	password string
}

const (
	authorizationKey = "Authorization"
	BlockName        = "BasicGuard"
)

func (bg BasicGuard) OnAttach(app *r.App) error {
	err := app.Decode("config", &bg.config)
	if err != nil {
		return err
	}
	return nil
}

func (bg BasicGuard) BasicAuthGuard(header http.Header) error {
	if bg.config.username == "" || bg.config.password == "" {
		return errors.New("BasicGuard: requires you to specify username & password inside" +
			"[basicguard] object of your config")
	}

	aHeader := header.Get(authorizationKey)
	if !strings.Contains(aHeader, "Basic") {
		return errors.New("This request doesn't have basic authorization key")
	}

	aHeader = strings.Replace(aHeader, "Basic ", "", 1)
	decoded, err := base64.StdEncoding.DecodeString(aHeader)
	if err != nil {
		return err
	}

	decData := strings.Split(string(decoded), ":")
	if len(decData) == 0 || len(decData) < 2 {
		return errors.New("Malformed Basic header.")
	}

	if decData[0] != bg.config.username || bg.config.password != decData[1] {
		return errors.New("Authorization failure!")
	}

	return nil
}

func init() {
	r.Attach(BlockName, BasicGuard{})
}

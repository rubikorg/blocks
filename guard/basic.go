package guard

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/rubikorg/rubik"
)

type BasicGuard struct{}

const (
	authorizationKey = "Authorization"
)

func (bg BasicGuard) Authorize(app *rubik.App, header http.Header) error {
	var config map[string]string

	err := app.Decode("basicguard", &config)
	if err != nil {
		return err
	}
	if config["username"] == "" || config["password"] == "" {
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

	if decData[0] != config["username"] || config["password"] != decData[1] {
		return errors.New("Authorization failure!")
	}

	return nil
}

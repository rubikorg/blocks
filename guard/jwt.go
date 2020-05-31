package guard

import (
	"github.com/rubikorg/rubik"
)

// JWT middleware
func JWT(a func(string) interface{}) rubik.Controller {
	return func(req *rubik.Request) {
		return a(req.Raw.Header.Get("Authorization"))
	}
}

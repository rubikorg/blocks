package healthcheck

import (
	"strings"

	"github.com/rubikorg/rubik"
	r "github.com/rubikorg/rubik"
)

// BlockName is the name of this block
const BlockName = "HealthCheck"

// BlockHealthCheck creates /health route for you and is used
// on service health checkers like kubernetes etc ..
type BlockHealthCheck struct {
	customPath string
}

var hcRoute = r.Route{
	Path: "/health",
	Controller: func(req *rubik.Request) {
		return req.Respond("ok")
	},
}

// OnAttach implementation for healthcheck block
func (hc BlockHealthCheck) OnAttach(app *r.App) error {
	conf := app.Config("healthcheck")
	if conf == nil {
		r.UseRoute(hcRoute)
	} else {
		c, ok := conf.(map[string]interface{})
		if ok && c["path"].(string) != "" {
			p := c["path"].(string)
			if !strings.HasPrefix(c["path"].(string), "/") {
				p = "/" + c["path"].(string)
			}
			hcRoute.Path = p
			r.UseRoute(hcRoute)
		}
	}
	return nil
}

func init() {
	r.Attach(BlockName, BlockHealthCheck{})
}

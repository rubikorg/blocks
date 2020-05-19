package healthcheck

import (
	"strings"

	r "github.com/rubikorg/rubik"
)

const BlockName = "HealthCheck"

type BlockHealthCheck struct {
	customPath string
}

var hcRoute = r.Route{
	Path:       "/health",
	Controller: hcCtl,
}

func hcCtl(en interface{}) r.ByteResponse {
	return r.Success("ok", r.Type.Text)
}

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

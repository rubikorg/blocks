package cors

import (
	"net/http"
	"strings"

	"github.com/rubikorg/rubik"
)

const (
	// BlockName is this Rubik block's name
	BlockName = "CORS"
)

// Options is used to parse your Rubik cors config
type Options struct {
	AllowedOrigins []string `json:"origins"`
	AllowedMethods []string `json:"methods"`

	AllowAllOrigins  bool `json:"allOrigins"`
	AllowAllMethods  bool `json:"allMethods"`
	AllowCredentials bool `json:"allowCredentials"`
}

// BlockCors defines your cors options for your server
type BlockCors struct {
	opts *Options
}

// OnAttach implementation for cors
func (c BlockCors) OnAttach(app *rubik.App) error {
	err := app.Decode("cors", c.opts)
	if err != nil {
		return err
	}
	return nil
}

// MW is the main cors middleware
func (c BlockCors) MW(extraOpts ...Options) rubik.Middleware {
	return func(req rubik.Request) interface{} {
		opts := c.opts
		if len(extraOpts) > 0 {
			opts = &extraOpts[0]
		}
		setHeadersFromOpts(opts, req)
		return nil
	}
}

func setHeadersFromOpts(opts *Options, req rubik.Request) {
	h := req.ResponseHeader
	r := req.Raw
	// preflight request
	if r.Method == http.MethodOptions && h.Get("Access-Control-Request-Method") != "" {
		h.Add("Vary", "Origin")
		h.Add("Vary", "Access-Control-Request-Method")
		h.Add("Vary", "Access-Control-Request-Headers")
	}

	if opts.AllowAllOrigins {
		h.Set("Access-Control-Allow-Origin", "*")
	} else {
		origins := strings.Join(opts.AllowedOrigins, ", ")
		h.Set("Access-Control-Allow-Origin", origins)
	}

	if opts.AllowAllMethods {
		methods := []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodOptions,
			http.MethodPatch,
		}
		h.Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
	} else {
		h.Set("Access-Control-Allow-Methods", strings.Join(opts.AllowedMethods, ", "))
	}

	if opts.AllowCredentials {
		h.Set("Access-Control-Allow-Credentials", "true")
	}

}

func init() {
	rubik.Attach(BlockName,
		BlockCors{&Options{
			AllowedMethods: []string{http.MethodGet, http.MethodOptions},
		}})
}

package swagger

import (
	"net/http"
	"strings"

	r "github.com/rubikorg/rubik"
)

// BlockName is the name of this rubik block
const BlockName = "Swagger"

var html = `
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger | Rubik</title>
    <link rel="stylesheet" type="text/css" href="http://localhost:5000/static/swagger-ui.css" >
    <link rel="icon" type="image/png" href="https://rubikorg.github.io/img/icon.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://rubikorg.github.io/img/icon.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="http://localhost:5000/static/swagger-ui-bundle.js"> </script>
    <script src="http://localhost:5000/static/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: '/rubik/docs/swagger.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region

      window.ui = ui
    }
  </script>
  </body>
</html>
`

var response = swagResponse{
	Swagger: "2.0",
	Info:    &info{},
	Tags:    []swagTag{},
	Paths:   make(swagPath),
	Schemes: []string{"http", "https"},
}

var swaggerRouter = r.Create("/rubik")

var htmlRoute = r.Route{
	Path:        "/docs",
	Description: "Serves the HTML body for Rubik Swagger Documentation",
}

var jsonRoute = r.Route{
	Path:        "/docs/swagger.json",
	Description: "Serves the RouteTree of Rubik as Swagger JSON",
}

// BlockSwagger is the swagger code block for Rubik server
type BlockSwagger struct{}

// OnAttach implementation of swagger
func (bs BlockSwagger) OnAttach(app *r.App) error {
	err := app.Decode("swagger", response.Info)
	if err != nil {
		return err
	}

	response.Host = app.CurrentURL
	insertSwaggerTags(app.RouterList)
	insertPaths(app.RouteTree.Routes)
	return nil
}

func (bs BlockSwagger) serve(en interface{}) r.ByteResponse {
	return r.Success(response, r.Type.JSON)
}

func (bs BlockSwagger) servePage(en interface{}) r.ByteResponse {
	return r.Success(html, r.Type.HTML)
}

func insertPaths(ri []r.RouteInfo) {
	for _, info := range ri {
		method := "get"
		belongsTo := "index"
		if info.BelongsTo != "" {
			belongsTo = info.BelongsTo
		}

		if info.Method != "" {
			method = strings.ToLower(info.Method)
		}

		responses := make(map[int]swagRespDecl)
		for status, respType := range info.Responses {
			resp := swagRespDecl{
				Description: http.StatusText(status),
				Schema: map[string]string{
					"type": respType,
				},
			}
			responses[status] = resp
		}

		pathInfo := map[string]swagPathInfo{
			method: {
				Tags:       []string{belongsTo},
				Summary:    info.Description,
				Parameters: []swagParams{},
				Produces:   []string{"application/json"},
				Responses:  responses,
			},
		}
		response.Paths[info.Path] = pathInfo
	}
}

func insertSwaggerTags(rl map[string]string) {
	for k, v := range rl {
		name := k
		if k == "" {
			name = "index"
		}
		t := swagTag{
			Name:        name,
			Description: v,
		}
		response.Tags = append(response.Tags, t)
	}
}

func init() {
	swaggerRouter.Description = "Swagger Implementation of Rubik - [@codekidX](https://github.com/codekidX)"

	block := BlockSwagger{}

	jsonRoute.Controller = block.serve
	htmlRoute.Controller = block.servePage
	swaggerRouter.Add(jsonRoute)
	swaggerRouter.Add(htmlRoute)
	r.Use(swaggerRouter)

	r.AttachAfter(BlockName, block)
}

package swagger

import (
	"net/http"
	"strings"

	"github.com/rubikorg/rubik"
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

// BlockSwagger is the swagger code block for Rubik server
type BlockSwagger struct{}

var response = swagResponse{
	Swagger: "2.0",
	Info:    &info{},
	Tags:    []swagTag{},
	Paths:   make(swagPath),
	Schemes: []string{"http", "https"},
}

type swagTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type swagPath map[string]map[string]swagPathInfo

type swagRespDecl struct {
	Description string            `json:"description"`
	Schema      map[string]string `json:"schema"`
}
type swagPathInfo struct {
	Summary    string               `json:"summary"`
	Tags       []string             `json:"tags"`
	Parameters []swagParams         `json:"parameters"`
	Produces   []string             `json:"produces"`
	Responses  map[int]swagRespDecl `json:"responses"`
}

type swagParams struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
	Format      string `json:"format"`
}

type swagResponse struct {
	Info    *info     `json:"info"`
	Swagger string    `json:"swagger"`
	Host    string    `json:"host"`
	Tags    []swagTag `json:"tags"`
	Paths   swagPath  `json:"paths"`
	Schemes []string  `json:"schemes"`
}

// Info is the info block of swagger guideline response
type info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Terms       string `json:"termsOfService"`
}

type swaggerEn struct {
	rubik.Entity
	AppURL string
}

// OnAttach implementation of swagger
func (bs BlockSwagger) OnAttach(app *rubik.App) error {
	err := app.Decode("swagger", response.Info)
	if err != nil {
		return err
	}

	response.Host = app.CurrentURL
	insertSwaggerTags(app.RouterList)
	insertPaths(app.RouteTree.Routes)
	return nil
}

func (bs BlockSwagger) serve(en interface{}) rubik.ByteResponse {
	return rubik.Success(response, rubik.Type.JSON)
}

func (bs BlockSwagger) servePage(en interface{}) rubik.ByteResponse {
	return rubik.Success(html, rubik.Type.HTML)
}

func insertPaths(ri []rubik.RouteInfo) {
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
			method: swagPathInfo{
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

var swaggerRouter = rubik.Create("/rubik")

var htmlRoute = rubik.Route{
	Path:        "/docs",
	Description: "Servers the HTML body for Rubik Swagger Documentation",
}

var jsonRoute = rubik.Route{
	Path:        "/docs/swagger.json",
	Description: "Servers the RouteTree of Rubik as Swagger JSON",
}

func init() {
	swaggerRouter.Description = "Swagger Implementation of Rubik - [@codekidX](https://github.com/codekidX)"

	block := BlockSwagger{}

	jsonRoute.Controller = block.serve
	htmlRoute.Controller = block.servePage
	swaggerRouter.Add(jsonRoute)
	swaggerRouter.Add(htmlRoute)
	rubik.Use(swaggerRouter)

	rubik.AttachAfter(BlockName, block)
}

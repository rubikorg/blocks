package swagger

import (
	"fmt"

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
    <title>Rubik Swagger Block</title>
    <link rel="stylesheet" type="text/css" href="http://localhost:5000/static/swagger-ui.css" >
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
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
        url: '%s',
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
type BlockSwagger struct {
	response swagResponse
	jsonURL *string
}

type swagTag struct {
	Name string
	Description string
}

type swagPath map[string]swagPathInfo

type swagPathInfo struct {
	Summary string
}

type swagResponse struct {
	Info *info `json:"info"`
	Swagger string `json:"swagger"`
	Host string `json:"host"`
	Tags []swagTag `json:"tags"`
	Paths map[string]swagPath `json:"paths"`
}

// Info is the info block of swagger guideline response
type info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type swaggerEn struct {
	rubik.Entity
	AppURL string
}

// OnAttach implementation of swagger
func (bs BlockSwagger) OnAttach(app *rubik.App) error {
	err := app.Decode("swagger", bs.response.Info)
	if err != nil {
		return err
	}
	blockJSONURL = app.CurrentURL + "/rubik/docs/swagger.json"
	return nil
}

func (bs BlockSwagger) serve(en interface{}) rubik.ByteResponse {
	return rubik.Success(bs.response, rubik.Type.JSON)
}

func (bs BlockSwagger) servePage(en interface{}) rubik.ByteResponse {
	return rubik.Success(fmt.Sprintf(html, *bs.jsonURL), rubik.Type.HTML)
}

var swaggerRouter = rubik.Create("/rubik")

var htmlRoute = rubik.Route{
	Path: "/docs",
}

var jsonRoute = rubik.Route{
	Path: "/docs/swagger.json",
}

var blockJSONURL = ""

func init() {
	block := BlockSwagger{
		response: swagResponse{
			Swagger: "2.0",
			Info: &info{},
		},
		jsonURL: &blockJSONURL,
	}
	rubik.AttachAfter(BlockName, block)


	jsonRoute.Controller = block.serve
	htmlRoute.Controller = block.servePage
	swaggerRouter.Add(jsonRoute)
	swaggerRouter.Add(htmlRoute)
	rubik.Use(swaggerRouter)

}
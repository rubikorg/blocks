package apigen

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	tsGen "github.com/rubikorg/blocks/apigen/ts/templates"
	r "github.com/rubikorg/rubik"
)

type config struct {
	outFolder string
}

// TSExtBlock generates typescript client code from Rubik routes
type TSExtBlock struct {
	conf config
}

// OnPlug satisfies the rubik.ExtentionBlock interface
func (TSExtBlock) OnPlug(app *r.App) error {
	var templateData = make(map[string]*tsGen.TypescriptTemplate)

	for router := range app.RouteTree.RouterList {
		var name string
		if router == "" {
			name = "Index"
		} else {
			name = capitalize(strings.ReplaceAll(router, "/", ""))
		}
		templateData[name] = &tsGen.TypescriptTemplate{
			RouterName: name,
			Routes:     []tsGen.TsRoute{},
		}
	}

	for _, route := range app.RouteTree.Routes {
		var tsRoute tsGen.TsRoute
		var routerName string

		if route.BelongsTo == "" {
			routerName = "Index"
		} else {
			routerName = capitalize(strings.ReplaceAll(route.BelongsTo, "/", ""))
		}

		target := templateData[routerName]

		if route.Entity != nil {
			values := reflect.ValueOf(route.Entity)
			tsRoute.EntityName = capitalize(values.Type().Name())
			tsRoute.Path = route.Path
			tsRoute.FullPath = route.FullPath
			if route.Path == "/" {
				tsRoute.Name = "root"
			} else {
				tsRoute.Name = uncapitalize(
					strings.ReplaceAll(
						handleDotRoutePath(route.Path), "/", ""))
			}

			if route.Method == "" {
				tsRoute.Method = http.MethodGet
			} else {
				tsRoute.Method = route.Method
			}
			num := values.NumField()

			// add entity data pairs
			for i := 0; i < num; i++ {
				field := values.Type().Field(i)

				tag := field.Tag.Get("rubik")
				key, medium := getRequestField(field.Name, tag)
				typ := getTsTypeEquivalent(field.Type.Name())
				if typ == "-1" {
					continue
				}

				pair := tsGen.Pair{
					Key:        key,
					IsOptional: true,
					Type:       typ,
				}
				switch medium {
				case "body":
					tsRoute.Body = append(tsRoute.Body, pair)
					break
				case "query":
					tsRoute.Query = append(tsRoute.Query, pair)
					break
				case "param":
					tsRoute.Param = append(tsRoute.Param, pair)
					break
				default:
					tsRoute.Query = append(tsRoute.Query, pair)
					break
				}
			}

			target.Routes = append(target.Routes, tsRoute)
		}
	}

	// TOOD: replace stdout with file buffer
	// all router files and it's APIs
	for file, data := range templateData {
		fmt.Println("File: ", file, "=>")
		tmpl, err := template.New("api_file").Parse(tsGen.APITemplate)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(os.Stdout, *data); err != nil {
			return err
		}
	}

	// env file
	tmpl, err := template.New("env_file").Parse(tsGen.ENVTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(os.Stdout, struct{ URL string }{app.CurrentURL}); err != nil {
		return err
	}

	// types file can be written as is

	return nil
}

// Name mentions the name of the extension to the use
func (TSExtBlock) Name() string {
	return "Typescript Client SDK Generator"
}

func uncapitalize(field string) string {
	r := []rune(field)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func capitalize(field string) string {
	r := []rune(field)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func handleDotRoutePath(path string) string {
	var r []rune
	var foundDot = false
	for _, p := range path {
		if p == '.' {
			foundDot = true
			continue
		}
		if foundDot {
			r = append(r, unicode.ToUpper(p))
			foundDot = false
		} else {
			r = append(r, p)
		}
	}

	return string(r)
}

// TODO: need to handle nested structs inside the type
func getTsTypeEquivalent(goType string, value ...reflect.Value) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float", "float32", "float64",
		"complex32", "complex64":
		return "number"
	case "bool":
		return "boolean"
	default:
		// ignore!
		return "-1"
	}
}

// getRequestField takes the original field name and rubik struct tag and
// returns (key_name, medium)
func getRequestField(ogName string, tag string) (string, string) {
	// TODO: handle optional entity fields
	key := uncapitalize(ogName)
	constTag := strings.ReplaceAll(tag, "!", "")
	if strings.Contains(constTag, "|") {
		splitted := strings.Split(constTag, "|")
		var medium = "query"
		if splitted[0] != "" {
			key = splitted[0]
		}

		if splitted[1] != "" {
			medium = splitted[1]
		}
		return key, medium
	}

	switch constTag {
	case "body", "query", "param":
		return key, constTag
	default:
		return key, "query"
	}
}

func init() {
	r.Plug(TSExtBlock{})
}

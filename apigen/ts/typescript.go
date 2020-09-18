package apigen

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"

	tsGen "github.com/rubikorg/blocks/apigen/ts/templates"
	r "github.com/rubikorg/rubik"
	"github.com/rubikorg/rubik/pkg"
)

type config struct {
	OutDir    string `json:"out_dir"`
	CompileJS bool   `json:"compile_js"`
}

// TSExtBlock generates typescript client code from Rubik routes
type TSExtBlock struct {
	conf config
}

var conf config

// OnPlug satisfies the rubik.ExtentionBlock interface
func (TSExtBlock) OnPlug(app *r.App) error {
	err := app.Decode("sdk_ts", &conf)
	if err != nil {
		return err
	}

	// if no out dir specifies rubik workspace
	if conf.OutDir == "" {
		conf.OutDir = filepath.Join("..", "..", "apigen", "ts")
	}

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
					replaceURLPathAsName(route.Path))
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

	outDir := conf.OutDir

	// if you want to compile to js
	if conf.CompileJS {
		// check if tsc is installed globally
		if _, err := exec.LookPath("tsc"); err != nil {
			return errors.New("Please install `tsc` as a global executable using: `npm i -g tsc`")
		}

		// create a outDir inside rubik cache
		outDir = filepath.Join(pkg.MakeAndGetCacheDirPath(), "apigen_ts")
		if f, _ := os.Stat(outDir); f == nil {
			os.MkdirAll(outDir, 0755)
		}
	}

	if f, _ := os.Stat(conf.OutDir); f == nil {
		os.MkdirAll(conf.OutDir, 0755)
	}

	// all router files and it's APIs
	var buf bytes.Buffer
	for file, data := range templateData {
		// if there is no routes in this router continue
		if len(data.Routes) == 0 {
			continue
		}

		tmpl, err := template.New("api_file").Parse(tsGen.APITemplate)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(&buf, *data); err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s-route.ts", strings.ToLower(file))
		filePath := filepath.Join(outDir, fileName)
		err = ioutil.WriteFile(filePath, buf.Bytes(), 0755)
		if err != nil {
			return err
		}

		buf.Reset()
	}

	// env file
	tmpl, err := template.New("env_file").Parse(tsGen.ENVTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buf, struct{ URL string }{app.CurrentURL}); err != nil {
		return err
	}

	// env file
	err = ioutil.WriteFile(filepath.Join(outDir, "rubik-env.ts"), buf.Bytes(), 0755)

	// all the files which does not need template data to be passed
	for name, tplData := range tsGen.TSFileMap {
		err = ioutil.WriteFile(filepath.Join(outDir, name), []byte(tplData), 0755)
		if err != nil {
			return err
		}
	}

	if conf.CompileJS {
		fmt.Println("Compiling Typescript project to Javascript...")
		cmd := exec.Command("tsc", fmt.Sprintf("%s/*.ts", outDir), "--outDir", conf.OutDir)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()

		fmt.Println("Cleaning up...")
		if err := os.RemoveAll(outDir); err != nil {
			return err
		}
	}

	fmt.Printf(`
Generated HTTP Request for your corresponding Rubik service:

path: %s
dependencies:
"axios"`, conf.OutDir)

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

func replaceURLPathAsName(path string) string {
	var r []rune
	var foundDot = false
	var foundColon = false
	for _, p := range path {
		if p == '/' {
			continue
		}
		if p == '.' {
			foundDot = true
			continue
		}
		if p == ':' {
			foundColon = true
			continue
		}
		if foundDot {
			r = append(r, unicode.ToUpper(p))
			foundDot = false
		} else if foundColon {
			r = append(r, unicode.ToUpper(p))
			foundColon = false
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

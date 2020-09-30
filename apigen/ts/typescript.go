package apigen

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/rubikorg/blocks/apigen/sdkdist"
	tsGen "github.com/rubikorg/blocks/apigen/ts/templates"
	r "github.com/rubikorg/rubik"
	"github.com/rubikorg/rubik/pkg"
)

type config struct {
	OutDir    string `json:"out_dir"`
	CompileJS bool   `json:"compile_js"`
}

// TSGenPlugin generates typescript client code from Rubik routes
type TSGenPlugin struct{}

var conf config

// OnPlug satisfies the rubik.ExtentionBlock interface
func (TSGenPlugin) OnPlug(app *r.App) error {
	// if --args is empty
	if app.Args == "" {
		// default path
		conf.OutDir = filepath.Join("..", "..", "apigen", "ts")
	} else {
		// split the args into consumable map. format - out=./out_dir,
		if out := sdkdist.MapPluginArgs(app.Args)["out"]; out != "" {
			conf.OutDir = out
		} else {
			return errors.New("out option not specified in args. please add out=$DIR in --args")
		}
	}

	var templateData = sdkdist.TransformTree(app.RouteTree, getTsTypeEquivalent)

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
func (TSGenPlugin) Name() string {
	return "Typescript Client SDK Generator"
}

func (TSGenPlugin) RunID() string {
	return "apigen_ts"
}

// TODO: need to handle nested structs inside the type
func getTsTypeEquivalent(goType string, value reflect.StructField) string {
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

func init() {
	r.Plug(TSGenPlugin{})
}

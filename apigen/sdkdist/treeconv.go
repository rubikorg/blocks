package sdkdist

import (
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/rubikorg/rubik"
)

type TypeConverterFunc func(string, reflect.StructField) string

// Pair is the object key: type values for { interface }
// definitions in TS
type Pair struct {
	Key        string
	IsOptional bool
	Type       string
}

// TsRoute is the route information for single route
// in the list of rubik.Router
type TreeRoute struct {
	FullPath   string
	Path       string
	Name       string
	Method     string
	EntityName string
	Form       []Pair
	Body       []Pair
	Query      []Pair
	Param      []Pair
}

// APITemplate is the rubik.Router -> Api Class
// definition in TS
type APITemplate struct {
	RouterName string
	Routes     []TreeRoute
}

func TransformTree(
	routeTree rubik.RouteTree,
	typeConverterFunc TypeConverterFunc) map[string]*APITemplate {
	var templateData = make(map[string]*APITemplate)
	for router := range routeTree.RouterList {
		var name string
		if router == "" {
			name = "Index"
		} else {
			name = capitalize(strings.ReplaceAll(router, "/", ""))
		}
		templateData[name] = &APITemplate{
			RouterName: name,
			Routes:     []TreeRoute{},
		}
	}

	for _, route := range routeTree.Routes {
		var tsRoute TreeRoute
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
				typ := typeConverterFunc(field.Type.Name(), field)
				if typ == "-1" {
					continue
				}

				pair := Pair{
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
	return templateData
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

func MapPluginArgs(args string) map[string]string {
	var a = make(map[string]string)
	var splitted []string
	if strings.Contains(args, ",") {
		splitted = strings.Split(args, ",")
		for _, s := range splitted {
			if strings.Contains(s, "=") {
				opt := strings.Split(strings.Trim(s, " "), "=")
				a[opt[0]] = opt[1]
			}
		}
	} else {
		if strings.Contains(args, "=") {
			opt := strings.Split(args, "=")
			a[opt[0]] = opt[1]
		}
	}
	return a
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

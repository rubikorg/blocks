package templates

// TODO: add this to rubik-api-helper.ts

// if (Object.Keys({{ $route.EntityName }}.query).length > 0) {
// 			reqUrl += '?' + encode()
// 		}

// {{ if eq $route.Method "GET" }}
// 		return getRequest();
// 		{{ else }} {{ if eq $route.Method "POST" }}
// 		return postRequest(env.url + {{ $route.Path }}, opts);
// 		{{ end }}
// 	{{ end }}

// interface ApiOptions {
// 	headers: [key:string]:string,
// 	auth: AuthOptions
// }

// Pair is the object key: type values for { interface }
// definitions in TS
type Pair struct {
	Key        string
	IsOptional bool
	Type       string
}

// TsRoute is the route information for single route
// in the list of rubik.Router
type TsRoute struct {
	FullPath   string
	Path       string
	Name       string
	Method     string
	EntityName string
	Body       []Pair
	Query      []Pair
	Param      []Pair
}

// TypescriptTemplate is the rubik.Router -> Api Class
// definition in TS
type TypescriptTemplate struct {
	RouterName string
	Routes     []TsRoute
}

// EnvFile is the settings and environment related
// values provider for Rubik SDK files
type EnvFile struct {
	URL string
}

// TSFileMap is the map of static files and corresponding templates
// which are to be written with no template variables, these files
// tend to be helper files or type definitions for TS lang
var TSFileMap = map[string]string{
	"rubik-env.ts":        ENVTemplate,
	"types.ts":            TypesTemplate,
	"rubik-api-helper.ts": APIHelperTemplate,
}

// APITemplate is the constant router file which defines a single
// router in TS which corresponds to rubik.Router
const APITemplate = `import { doRequest } from './rubik-api-helper.ts';
import { env } from './rubik-env.ts';
import { ApiOptions } from './types.ts';
import { encode } from 'qs'; 
{{range $route := .Routes }}
interface {{ $route.EntityName }} {
	{{- if not $route.Body }}
	body?: any,
	{{- else }}
	body: {
	{{- range $bodyElem := $route.Body }}
		{{ $bodyElem.Key }}{{- if $bodyElem.IsOptional }}?{{- end }}: {{ $bodyElem.Type }},
	{{- end }}
	},
	{{- end }}
	{{- if not $route.Query }}
	query?: any,
	{{- else }}
	query: {
	{{- range $queryElem := $route.Query }}
		{{ $queryElem.Key }}{{- if $queryElem.IsOptional }}?{{- end }}: {{ $queryElem.Type }},
	{{- end }}
	},
	{{- end }}
	{{- if not $route.Param }}
	param?: any,
	{{- else }}
	param: {
	{{- range $paramElem := $route.Param }}
		{{ $paramElem.Key }}: {{ $paramElem.Type }},
	{{- end }}
	}
	{{- end }}
}
{{end}}
// @class {{ .RouterName }}Api
export class {{ .RouterName }}Api {
{{ range $route := .Routes }}
	public static {{ $route.Name }}(payload: {{ $route.EntityName }}, opts: ApiOptions): Promise<any> {
		return doRequest('{{ $route.Method }}', env.url + {{ $route.FullPath }}, opts);
	}
{{ end }}
}
`

// ENVTemplate is the template for rubik-env.ts file
const ENVTemplate = `
export const env = {
	url: '{{ .URL }}',
};
`

// TypesTemplate includes all types required for TS SDK
const TypesTemplate = `
export interface AuthOptions {
	basic: {
		username: string,
		password: string
	},
	jwt: string
}
`

// APIHelperTemplate implements the HTTP methods for Rubik TS SDK
const APIHelperTemplate = `
import { ApiOptions } from './types.ts';

export function doRequest(method: string, opts?: AuthOptions): Promise<any> {
	return Promise.resolve({});
}
`

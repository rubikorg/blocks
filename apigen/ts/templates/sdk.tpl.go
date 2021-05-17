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

// EnvFile is the settings and environment related
// values provider for Rubik SDK files
type EnvFile struct {
	URL string
}

// TSFileMap is the map of static files and corresponding templates
// which are to be written with no template variables, these files
// tend to be helper files or type definitions for TS lang
var TSFileMap = map[string]string{
	"types.ts":            TypesTemplate,
	"rubik-api-helper.ts": APIHelperTemplate,
}

// APITemplate is the constant router file which defines a single
// router in TS which corresponds to rubik.Router
const APITemplate = `import { doRequest } from './rubik-api-helper';
import { env } from './rubik-env';
import { Entity } from './types';
import { AxiosResponse } from "axios";
{{range $route := .Routes }}
export interface {{ $route.EntityName }} extends Entity {
	{{- if not $route.Form }}
	form?: any,
	{{- else }}
	form: {
	{{- range $bodyElem := $route.Form }}
		{{ $bodyElem.Key }}{{- if $bodyElem.IsOptional }}?{{- end }}: {{ $bodyElem.Type }},
	{{- end }}
	},
	{{- end }}
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
	public static {{ $route.Name }}(entity: {{ $route.EntityName }}): Promise<AxiosResponse> {
		return doRequest<{{ $route.EntityName }}>('{{ $route.Method }}', env.url + "{{ $route.FullPath }}", entity);
	}
{{ end }}
}
`

// ENVTemplate is the template for rubik-env.ts file
const ENVTemplate = `export const env = {
	url: '{{ .URL }}',
};
`

// TypesTemplate includes all types required for TS SDK
const TypesTemplate = `export interface Entity {
	query?: any;
	body?: any;
	form?: any;
	param?: any;
	auth?: {
		basic: {
			username: string,
			password: string
		},
		jwt: string
	}
}
`

// APIHelperTemplate implements the HTTP methods for Rubik TS SDK
const APIHelperTemplate = `import { Entity } from "./types";
import axios from "axios";
import { encode } from "qs";

export async function doRequest<T extends Entity>(
	method: "GET" | "POST" | "PUT" | "DELETE",
	path: string,
	entity: T
): Promise<any> {
	switch (method) {
		case "GET":
			return getRequest(path, entity);
		case "POST":
			return postRequest(path, entity);
		case "PUT":
			return putRequest(path, entity);
		case "DELETE":
			return deleteRequest(path, entity);
		default:
			break;
	}
}

function getRequest<T extends Entity>(path: string, entity: T): Promise<any> {
	// TODO: we can improve this by staticly analyzing while creating template
	let finalPath = path;
	if (Object.keys(entity.param).length > 0) {
		// then we have some path params to be passed and replaced in url path
		Object.keys(entity.param).forEach((k) => {
			finalPath = finalPath.replace(":" + k, entity.param[k]);
		});
	}

	if (Object.keys(entity.query).length > 0) {
		const encodedQuery = encode(entity.query);
		finalPath += "?" + encodedQuery;
	}

	if (Object.keys(entity.body).length > 0) {
		console.log("[Rubik] Cannot embed body in a GET request");
	}

	return axios.get(finalPath);
}

function evalRequestData(path: string, entity: Entity): [string, any] {
	let finalPath = path;
	// TODO: we can improve this by staticly analyzing while creating template
	if (Object.keys(entity.param).length > 0) {
		// then we have some path params to be passed and replaced in url path
		Object.keys(entity.param).forEach((k) => {
			finalPath = finalPath.replace(":" + k, entity.param[k]);
		});
	}

	if (Object.keys(entity.query).length > 0) {
		const encodedQuery = encode(entity.query);
		finalPath += "?" + encodedQuery;
	}

	return [finalPath, entity.body];
}

function postRequest<T extends Entity>(path: string, entity: T): Promise<any> {
	const [finalPath, data] = evalRequestData(path, entity);
	return axios.post(finalPath, data);
}

function putRequest<T extends Entity>(path: string, entity: T): Promise<any> {
	const [finalPath, data] = evalRequestData(path, entity);
	return axios.put(finalPath, data);
}

function deleteRequest<T extends Entity>(
	path: string,
	entity: T
): Promise<any> {
	const [finalPath, data] = evalRequestData(path, entity);
	return axios.delete(finalPath, data);
}
`

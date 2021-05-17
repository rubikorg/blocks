package swagger

import r "github.com/rubikorg/rubik"

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
	r.Entity
	AppURL string
}

package template

import _ "embed"

//go:embed app.tmpl
var text string

func getTmpl() string {
	return text
}
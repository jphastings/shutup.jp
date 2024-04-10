package main

import (
	_ "embed"
	"html/template"
)

func lastPC(pcs []Postcard) Postcard {
	return pcs[len(pcs)-1]
}

func safeHTML(str string) template.HTML {
	return template.HTML(str)
}

//go:embed index.html.tmpl
var indexStr string
var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"safeHTML": safeHTML,
	"lastPC":   lastPC,
}).Parse(indexStr))

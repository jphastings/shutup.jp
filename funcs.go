package main

import (
	_ "embed"
	ht "html/template"
	"strings"
	tt "text/template"
	"time"

	"github.com/gohugoio/hugo/tpl"
)

func now() time.Time {
	return time.Now()
}

func lastPC(pcs []Postcard) Postcard {
	return pcs[len(pcs)-1]
}

func safeHTML(str string) ht.HTML {
	return ht.HTML(str)
}

func plainify(str string) string {
	return strings.TrimSpace(tpl.StripHTML(str))
}

var funcs = ht.FuncMap{
	"safeHTML": safeHTML,
	"lastPC":   lastPC,
	"now":      now,
	"plainify": plainify,
}

//go:embed index.html.tmpl
var indexStr string
var indexTmpl = ht.Must(ht.New("index").Funcs(funcs).Parse(indexStr))

//go:embed feed.xml.tmpl
var feedStr string
var feedTmpl = tt.Must(tt.New("feed").Funcs(funcs).Parse(feedStr))

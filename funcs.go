package main

import (
	_ "embed"
	"fmt"
	ht "html/template"
	"sort"
	"strings"
	tt "text/template"
	"time"

	"github.com/gohugoio/hugo/tpl"
	"github.com/jphastings/dotpostcard/formats"
	"github.com/jphastings/dotpostcard/types"
)

func now() time.Time {
	return time.Now()
}

func lastPC(pcs []types.Postcard) types.Postcard {
	return pcs[len(pcs)-1]
}

func htmlEscape(str string) string {
	return tt.HTMLEscapeString(str)
}

func plainify(str string) string {
	return strings.TrimSpace(tpl.StripHTML(str))
}

func altText(meta types.Metadata) string {
	alt, _ := formats.AltText(meta, "en")
	return alt
}

func annotate(at types.AnnotatedText) ht.HTML {
	sort.Sort(types.SortEndsLast(at.Annotations))

	// TODO: Handle annotations
	fmt.Println(at.Annotations)

	return ht.HTML(htmlEscape(at.Text))
}

var funcs = ht.FuncMap{
	"htmlEscape": htmlEscape,
	"lastPC":     lastPC,
	"now":        now,
	"plainify":   plainify,
	"altText":    altText,
	"annotate":   annotate,
}

//go:embed index.html.tmpl
var indexStr string
var indexTmpl = ht.Must(ht.New("index").Funcs(funcs).Parse(indexStr))

//go:embed feed.xml.tmpl
var feedStr string
var feedTmpl = tt.Must(tt.New("feed").Funcs(funcs).Parse(feedStr))

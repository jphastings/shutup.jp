package main

import (
	_ "embed"
	ht "html/template"
	"strings"
	tt "text/template"
	"time"

	"github.com/gohugoio/hugo/tpl"
	"github.com/jphastings/dotpostcard/formats"
	"github.com/jphastings/dotpostcard/types"
	"github.com/jphastings/shutup.jp/build/mapping"
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
	return ht.HTML(at.HTML())
}

//go:embed world.svg
var worldMapSVG string

func worldMap(cs mapping.Countries, pcs []types.Postcard) ht.HTML {
	collectedMap := worldMapSVG

	for _, cc := range cs.Codes() {
		collectedMap = strings.ReplaceAll(collectedMap, `id='`+cc+`'`, `id='`+cc+`' class='collected'`)
	}

	var allPoints string
	for _, pc := range pcs {
		x, y := ProjectRobinson(*pc.Meta.Location.Latitude, *pc.Meta.Location.Longitude)
		allPoints += svgPointer(x, y, pc.Meta.Location.CountryCode)
	}

	collectedMap = strings.Replace(collectedMap, "id='points'>", "id='points'>"+allPoints, 1)

	return ht.HTML(collectedMap)
}

var funcs = ht.FuncMap{
	"htmlEscape": htmlEscape,
	"lastPC":     lastPC,
	"now":        now,
	"plainify":   plainify,
	"altText":    altText,
	"annotate":   annotate,
	"worldMap":   worldMap,
}

//go:embed index.html.tmpl
var indexStr string
var indexTmpl = ht.Must(ht.New("index").Funcs(funcs).Parse(indexStr))

//go:embed feed.xml.tmpl
var feedStr string
var feedTmpl = tt.Must(tt.New("feed").Funcs(funcs).Parse(feedStr))

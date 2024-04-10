package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
)

//go:embed index.html.tmpl
var indexStr string
var indexTmpl = template.Must(template.New("index").Funcs(template.FuncMap{
	"safeHTML": func(str string) template.HTML {
		return template.HTML(str)
	},
}).Parse(indexStr))

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files, err := filepath.Glob("postcards/*.json")
	check(err)

	var pcs []Postcard
	toCopy := []string{"static/postcard.css", "static/shutup.css"}

	for _, file := range files[:3] {
		var pc Postcard
		f, err := os.Open(file)
		check(err)
		check(json.NewDecoder(f).Decode(&pc))

		pc.Name = filepath.Base(file)[:len(filepath.Base(file))-5]
		pcs = append(pcs, pc)
		toCopy = append(toCopy, "postcards/"+pc.Name+".webp")
	}

	indexF, err := os.OpenFile("dist/index.html", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)

	sort.Sort(BySentOn(pcs))
	check(indexTmpl.Execute(indexF, pcs))

	for _, tc := range toCopy {
		src, err := os.Open(tc)
		check(err)

		dst, err := os.OpenFile("dist/"+filepath.Base(tc), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		check(err)

		_, err = io.Copy(dst, src)
		check(err)
	}
}

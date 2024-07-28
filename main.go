package main

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/jphastings/dotpostcard/formats"
	"github.com/jphastings/dotpostcard/formats/component"
	"github.com/jphastings/dotpostcard/formats/web"
	"github.com/jphastings/dotpostcard/types"
)

const outDir = "dist/"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type tmplVars struct {
	Postcards []types.Postcard
	Sizes     map[string]int64
}

func main() {
	files, err := filepath.Glob("postcards/*.postcard.webp")
	check(err)

	vars := tmplVars{
		Sizes: make(map[string]int64),
	}
	toCopy := []string{
		"static/postcard.css",
		"static/shutup.css",
		"static/bg-light.png",
		"static/bg-dark.png",
		"static/og-image.jpg",
	}
	check(os.MkdirAll(outDir, 0755))

	for _, file := range files {
		f, err := os.Open(file)
		check(err)
		defer f.Close()

		b := web.BundleFromReader(f, file)
		pc, err := b.Decode(nil)
		check(err)

		pcImgName := "postcards/" + pc.Name + ".postcard.webp"
		vars.Postcards = append(vars.Postcards, pc)
		toCopy = append(toCopy, pcImgName)

		// Make front & back covers
		for _, fw := range component.Codec().Encode(pc, &formats.EncodeOptions{MaxDimension: 800}) {
			fname, err := fw.WriteFile(outDir, false)
			if !errors.Is(err, os.ErrExist) {
				check(err)
			}
			fs, err := os.Stat(path.Join(outDir, fname))
			check(err)
			vars.Sizes[fname] = fs.Size()
		}
	}
	sort.Sort(types.BySentOn(vars.Postcards))

	indexF, err := os.OpenFile(path.Join(outDir, "index.html"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)
	check(indexTmpl.Execute(indexF, vars))

	feedF, err := os.OpenFile(path.Join(outDir, "feed.xml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)
	check(feedTmpl.Execute(feedF, vars))

	for _, tc := range toCopy {
		src, err := os.Open(tc)
		check(err)

		dst, err := os.OpenFile(path.Join(outDir, filepath.Base(tc)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		check(err)

		_, err = io.Copy(dst, src)
		check(err)
	}
}

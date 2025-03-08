package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/jphastings/dotpostcard/formats"
	"github.com/jphastings/dotpostcard/formats/component"
	"github.com/jphastings/dotpostcard/formats/web"
	"github.com/jphastings/dotpostcard/types"
	"github.com/jphastings/shutup.jp/build/mapping"
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
	Countries mapping.Countries
}

func main() {
	files, err := filepath.Glob("postcards/*.postcard.webp")
	check(err)

	vars := tmplVars{
		Sizes:     make(map[string]int64),
		Countries: make(mapping.Countries),
	}
	toCopy := []string{
		"static/postcard.css",
		"static/shutup.css",
		"static/bg-light.png",
		"static/bg-dark.png",
		"static/og-image.jpg",
	}
	check(os.MkdirAll(outDir, 0755))

	fmt.Println("Finding postcards & extracting front/back images")
	for _, file := range files {
		f, err := os.Open(file)
		check(err)
		defer f.Close()

		b := web.BundleFromReader(f, file)
		pc, err := b.Decode(formats.DecodeOptions{})
		check(err)

		pcImgName := "postcards/" + pc.Name + ".postcard.webp"
		vars.Postcards = append(vars.Postcards, pc)
		toCopy = append(toCopy, pcImgName)
		vars.Countries.Add(pc.Meta.Location)

		// Make front & back covers
		fws, err := component.Codec().Encode(pc, &formats.EncodeOptions{MaxDimension: 800})
		check(err)

		for _, fw := range fws {
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

	fmt.Println("Creating index.html")
	indexF, err := os.OpenFile(path.Join(outDir, "index.html"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)
	check(indexTmpl.Execute(indexF, vars))

	fmt.Println("Creating feed.xml")
	feedF, err := os.OpenFile(path.Join(outDir, "feed.xml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)
	check(feedTmpl.Execute(feedF, vars))

	fmt.Printf("Copying files for distribution:\n")
	for _, tc := range toCopy {
		fmt.Printf("  %s\n", tc)
		src, err := os.Open(tc)
		check(err)

		dst, err := os.OpenFile(path.Join(outDir, filepath.Base(tc)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		check(err)

		_, err = io.Copy(dst, src)
		check(err)
	}
}

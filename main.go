package main

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/jphastings/dotpostcard/formats/metadata"
	"github.com/jphastings/dotpostcard/types"
)

const outDir = "dist/"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files, err := filepath.Glob("postcards/*-meta.json")
	check(err)

	var pcs []types.Postcard
	toCopy := []string{
		"static/postcard.css",
		"static/shutup.css",
		"static/bg-light.png",
		"static/bg-dark.png",
		"static/og-image.jpg",
	}

	for _, file := range files {
		f, err := os.Open(file)
		check(err)
		defer f.Close()

		b, err := metadata.BundleFromFile(f, path.Dir(file))
		check(err)

		pc, err := b.Decode(nil)
		check(err)

		imgName := "postcards/" + pc.Name + ".postcard"
		// fs, err := os.Stat(imgName)
		// check(err)
		// pc.ImgSize = fs.Size()

		pcs = append(pcs, pc)
		toCopy = append(toCopy, imgName)
	}
	sort.Sort(types.BySentOn(pcs))

	check(os.MkdirAll("dist/", 0755))

	indexF, err := os.OpenFile(path.Join(outDir, "index.html"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)
	check(indexTmpl.Execute(indexF, pcs))

	feedF, err := os.OpenFile(path.Join(outDir, "feed.xml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	check(err)
	check(feedTmpl.Execute(feedF, pcs))

	for _, tc := range toCopy {
		src, err := os.Open(tc)
		check(err)

		dst, err := os.OpenFile(path.Join(outDir, filepath.Base(tc)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		check(err)

		_, err = io.Copy(dst, src)
		check(err)
	}
}

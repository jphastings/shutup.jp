package main

import (
	"github.com/twpayne/go-proj/v11"
)

func ProjectRobinson(latitude, longitude float64) (x, y float64) {
	pj, err := proj.NewCRSToCRS("EPSG:4326", "+proj=robin", nil)
	if err != nil {
		panic(err)
	}

	robinson, err := pj.Forward(proj.NewCoord(latitude, longitude, 0, 0))
	if err != nil {
		panic(err)
	}

	x = (robinson.X()/1.700583333052523e+07)/2 + 0.5
	y = 0.5 - (robinson.Y()/8.6251546651e+06)/2

	return
}

package main

import "fmt"

var (
	// These figures represent the SVG map's transform compared to a full robinson projection
	offsetX = -33.89079281
	offsetY = -15.91771089
	scaleX  = 1.020619693 * 2000
	scaleY  = 1.207910396 * 857
)

func svgPointer(x, y float64, cc string) string {
	return fmt.Sprintf(
		`<circle r="5" transform="translate(%.1f,%.1f)"/>`,
		scaleX*x+offsetX,
		scaleY*y+offsetY,
	)
}

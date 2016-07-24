package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/bdf"
)

func main() {
	font, err := bdf.Parse("fonts/7x13B.bdf")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Font loaded: %d chars\n", font.Chars)
	glyph := font.Glyphs[65]
	fmt.Printf("A: height: %d len: %d\n", glyph.Height, len(glyph.Bitmap))
	for _, row := range glyph.Bitmap {
		for _, b := range row {
			if b {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

package main

import (
	"fmt"
	"github.com/depili/go-rgb-led-matrix/bdf"
	"github.com/depili/go-rgb-led-matrix/matrix"
	"time"
)

func main() {
	font, err := bdf.Parse("fonts/7x13B.bdf")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Font loaded: %d chars\n", font.Chars)
	glyph := font.GetGlyph(rune(65))
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

	text := "Testi ÄÖ"
	bitmap := font.TextBitmap(text)
	for _, row := range bitmap {
		for _, b := range row {
			if b {
				fmt.Printf("X")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}

	m := matrix.Init("tcp://192.168.0.30:5555", 32, 128)
	var color [3]byte
	color[0] = 0
	color[1] = 255
	color[2] = 0
	m.Fill(color)
	m.Send()
	time.Sleep(500 * time.Millisecond)
	color[0] = 255
	color[2] = 255
	m.Fill(matrix.ColorWhite())
	m.Send()
	time.Sleep(500 * time.Millisecond)

	m.InitFlame()
	for i := 0; i < 1000; i++ {
		m.FlameFill()
		m.Send()
		time.Sleep(10 * time.Millisecond)
	}

	m.Close()
}

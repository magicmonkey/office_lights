// +build ignore

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const size = 120

func main() {
	createIcon("icons/ledstrip.png", "LED\nStrip", color.RGBA{255, 100, 100, 255})
	createIcon("icons/ledbar_rgbw.png", "LED Bar\nRGBW", color.RGBA{150, 150, 255, 255})
	createIcon("icons/ledbar_white.png", "LED Bar\nWhite", color.RGBA{255, 255, 255, 255})
	createIcon("icons/videolight.png", "Video\nLights", color.RGBA{255, 200, 100, 255})
}

func createIcon(filename, text string, col color.RGBA) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Background
	bg := color.RGBA{30, 30, 30, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Draw colored rectangle
	rectSize := 60
	rectX := (size - rectSize) / 2
	rectY := 20
	rect := image.Rect(rectX, rectY, rectX+rectSize, rectY+rectSize)
	draw.Draw(img, rect, &image.Uniform{col}, image.Point{}, draw.Src)

	// Draw text below rectangle
	drawText(img, text, size/2, 90, color.RGBA{200, 200, 200, 255})

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func drawText(img *image.RGBA, text string, x, y int, col color.Color) {
	// Split text by newlines
	lines := splitLines(text)

	for i, line := range lines {
		point := fixed.Point26_6{
			X: fixed.I(x - len(line)*7/2),
			Y: fixed.I(y + i*15),
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(line)
	}
}

func splitLines(text string) []string {
	var lines []string
	line := ""
	for _, c := range text {
		if c == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(c)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

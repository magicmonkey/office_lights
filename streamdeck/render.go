package streamdeck

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	buttonWidth  = 120
	buttonHeight = 120
	touchWidth   = 800
	touchHeight  = 100
	sectionWidth = 200 // touchWidth / 4
)

// updateButtons renders and updates all button images
func (s *StreamDeckUI) updateButtons() error {
	for i := 0; i < 8; i++ {
		img, err := s.renderButton(i)
		if err != nil {
			return err
		}
		if err := s.device.SetImage(uint8(i), img); err != nil {
			return err
		}
		s.buttonImages[i] = img
	}
	return nil
}

// renderButton creates an image for a specific button
func (s *StreamDeckUI) renderButton(index int) (image.Image, error) {
	// Top row (0-3): Mode selection buttons
	if index < 4 {
		return s.renderModeButton(index)
	}

	// Bottom row (4-7): Reserved/unused
	return s.renderBlankButton(), nil
}

// renderModeButton renders a mode selection button
func (s *StreamDeckUI) renderModeButton(index int) (image.Image, error) {
	mode := Mode(index)
	isActive := s.currentMode == mode

	// Try to load icon from file
	iconPath := filepath.Join("streamdeck", "icons", s.getModeIconFilename(mode))
	img, err := loadImage(iconPath)
	if err != nil {
		// If icon not found, create a simple text button
		return s.renderTextButton(mode.String(), isActive), nil
	}

	// Apply highlight if active
	if isActive {
		img = applyHighlight(img)
	} else {
		img = applyDim(img)
	}

	return img, nil
}

// renderTextButton creates a simple text-based button
func (s *StreamDeckUI) renderTextButton(text string, isActive bool) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, buttonWidth, buttonHeight))

	// Background color
	bg := color.RGBA{30, 30, 30, 255}
	if isActive {
		bg = color.RGBA{60, 100, 180, 255}
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Draw text
	col := color.RGBA{200, 200, 200, 255}
	if isActive {
		col = color.RGBA{255, 255, 255, 255}
	}
	drawCenteredText(img, text, col)

	return img
}

// renderBlankButton creates a blank (black) button
func (s *StreamDeckUI) renderBlankButton() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, buttonWidth, buttonHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return img
}

// getModeIconFilename returns the filename for a mode's icon
func (s *StreamDeckUI) getModeIconFilename(mode Mode) string {
	switch mode {
	case ModeLEDStrip:
		return "ledstrip.png"
	case ModeLEDBarRGBW:
		return "ledbar_rgbw.png"
	case ModeLEDBarWhite:
		return "ledbar_white.png"
	case ModeVideoLights:
		return "videolight.png"
	default:
		return "unknown.png"
	}
}

// updateTouchscreen renders and updates the touchscreen display
// Note: The github.com/muesli/streamdeck library does not support the Stream Deck+ touchscreen.
// This method is kept for future implementation with a compatible library.
func (s *StreamDeckUI) updateTouchscreen() error {
	// TODO: Implement touchscreen support when a compatible library is available
	// For now, we just render the image but don't send it to the device
	s.touchImage = s.renderTouchscreen()
	return nil
}

// renderTouchscreen creates the full touchscreen image
func (s *StreamDeckUI) renderTouchscreen() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, touchWidth, touchHeight))

	// Background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{20, 20, 20, 255}}, image.Point{}, draw.Src)

	// Get section data for current mode
	sections := s.getSectionData()

	// Render each section
	for i := 0; i < 4; i++ {
		s.renderSection(img, i, sections[i])
	}

	return img
}

// renderSection renders one section of the touchscreen
func (s *StreamDeckUI) renderSection(img *image.RGBA, index int, data SectionData) {
	x := index * sectionWidth
	bounds := image.Rect(x, 0, x+sectionWidth, touchHeight)

	if !data.Active {
		// Inactive section - dark background
		draw.Draw(img, bounds, &image.Uniform{color.RGBA{10, 10, 10, 255}}, image.Point{x, 0}, draw.Src)
		return
	}

	// Active section - draw background
	bgColor := color.RGBA{40, 40, 40, 255}
	draw.Draw(img, bounds, &image.Uniform{bgColor}, image.Point{x, 0}, draw.Src)

	// Draw border
	borderColor := color.RGBA{80, 80, 80, 255}
	drawVerticalLine(img, x+sectionWidth-1, 0, touchHeight, borderColor)

	// Draw label at top
	labelY := 15
	drawTextAt(img, data.Label, x+sectionWidth/2, labelY, color.RGBA{150, 150, 150, 255}, true)

	// Draw value in center
	valueY := 50
	valueText := formatValue(data.Value)
	drawTextAt(img, valueText, x+sectionWidth/2, valueY, color.RGBA{255, 255, 255, 255}, true)

	// Draw progress bar at bottom
	barY := 80
	barHeight := 10
	barWidth := sectionWidth - 20
	barX := x + 10
	s.drawProgressBar(img, barX, barY, barWidth, barHeight, data.Value)
}

// drawProgressBar draws a horizontal progress bar
func (s *StreamDeckUI) drawProgressBar(img *image.RGBA, x, y, width, height, value int) {
	// Background
	bgRect := image.Rect(x, y, x+width, y+height)
	draw.Draw(img, bgRect, &image.Uniform{color.RGBA{60, 60, 60, 255}}, image.Point{x, y}, draw.Src)

	// Fill
	fillWidth := (value * width) / 255
	if fillWidth > 0 {
		fillRect := image.Rect(x, y, x+fillWidth, y+height)
		fillColor := color.RGBA{100, 180, 255, 255}
		draw.Draw(img, fillRect, &image.Uniform{fillColor}, image.Point{x, y}, draw.Src)
	}
}

// Helper functions

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func applyHighlight(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, img, image.Point{}, draw.Src)

	// Apply brightness increase
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := result.At(x, y).RGBA()
			r = min(r+0x2000, 0xffff)
			g = min(g+0x2000, 0xffff)
			b = min(b+0x2000, 0xffff)
			result.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return result
}

func applyDim(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, img, image.Point{}, draw.Src)

	// Apply brightness decrease
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := result.At(x, y).RGBA()
			r = r / 2
			g = g / 2
			b = b / 2
			result.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	return result
}

func drawCenteredText(img *image.RGBA, text string, col color.Color) {
	bounds := img.Bounds()
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2
	drawTextAt(img, text, centerX, centerY, col, true)
}

func drawTextAt(img *image.RGBA, text string, x, y int, col color.Color, centered bool) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	if centered {
		// Rough centering (basicfont has ~7 pixel width per char)
		textWidth := len(text) * 7
		point.X -= fixed.I(textWidth / 2)
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

func drawVerticalLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
	for y := y1; y < y2; y++ {
		img.Set(x, y, col)
	}
}

func formatValue(value int) string {
	if value < 100 {
		return string(rune('0' + value/10)) + string(rune('0' + value%10))
	}
	return string(rune('0'+value/100)) + string(rune('0'+(value%100)/10)) + string(rune('0'+value%10))
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

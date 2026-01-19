package streamdeck

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	sdlib "rafaelmartins.com/p/streamdeck"
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
	keyIDs := []sdlib.KeyID{
		sdlib.KEY_1, sdlib.KEY_2, sdlib.KEY_3, sdlib.KEY_4,
		sdlib.KEY_5, sdlib.KEY_6, sdlib.KEY_7, sdlib.KEY_8,
	}

	for i, keyID := range keyIDs {
		img, err := s.renderButton(i)
		if err != nil {
			return err
		}
		if err := s.device.SetKeyImage(keyID, img); err != nil {
			log.Printf("Warning: Failed to set key %d image: %v", i+1, err)
			// Continue even if one button fails
			continue
		}
		s.buttonImages[i] = img
	}
	return nil
}

// renderButton creates an image for a specific button
func (s *StreamDeckUI) renderButton(index int) (image.Image, error) {
	// Top row (0-3): Tab selection buttons
	if index < 4 {
		return s.renderTabButton(index)
	}

	// Second row (4-7): Tab-specific buttons
	switch s.currentTab {
	case TabLightControl:
		return s.renderModeButton(index - 4)
	case TabScenes:
		return s.renderSceneButton(index - 4)
	default:
		// Future tabs: show blank buttons
		return s.renderBlankButton(), nil
	}
}

// renderTabButton renders a tab selection button
func (s *StreamDeckUI) renderTabButton(index int) (image.Image, error) {
	tab := Tab(index)
	isActive := s.currentTab == tab

	// Try to load tab icon from file
	iconPath := filepath.Join("streamdeck", "icons", s.getTabIconFilename(tab))
	img, err := loadImage(iconPath)
	if err != nil {
		// If icon not found, create a simple text button
		return s.renderTextButton(tab.String(), isActive), nil
	}

	// Apply highlight if active
	if isActive {
		img = applyHighlight(img)
	} else {
		img = applyDim(img)
	}

	return img, nil
}

// getTabIconFilename returns the filename for a tab's icon
func (s *StreamDeckUI) getTabIconFilename(tab Tab) string {
	switch tab {
	case TabLightControl:
		return "tab_lights.png"
	case TabScenes:
		return "tab_scenes.png"
	case TabFuture3:
		return "tab_3.png"
	case TabFuture4:
		return "tab_4.png"
	default:
		return "unknown.png"
	}
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
func (s *StreamDeckUI) updateTouchscreen() error {
	// Check if device supports touchscreen
	if !s.device.GetTouchStripSupported() {
		// Device doesn't have a touch strip, skip
		return nil
	}

	img := s.renderTouchscreen()
	if err := s.device.SetTouchStripImage(img); err != nil {
		// Log but don't fail - touchscreen might not be available
		log.Printf("Warning: Failed to set touchscreen image: %v", err)
		return nil
	}
	s.touchImage = img
	return nil
}

// renderSceneButton renders a scene slot button
func (s *StreamDeckUI) renderSceneButton(index int) (image.Image, error) {
	exists, _ := s.storage.SceneExists(index)
	label := fmt.Sprintf("Scene %d", index+1)

	// Different styling for saved vs empty scenes
	return s.renderTextButton(label, exists), nil
}

// renderTouchscreen creates the full touchscreen image
func (s *StreamDeckUI) renderTouchscreen() image.Image {
	switch s.currentTab {
	case TabLightControl:
		return s.renderLightControlTouchscreen()
	case TabScenes:
		return s.renderScenesTouchscreen()
	default:
		return s.renderPlaceholderTouchscreen()
	}
}

// renderScenesTouchscreen renders the touchscreen for Tab 2 (Scenes)
func (s *StreamDeckUI) renderScenesTouchscreen() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, touchWidth, touchHeight))

	// Background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{20, 20, 20, 255}}, image.Point{}, draw.Src)

	// Render each scene section
	for i := 0; i < 4; i++ {
		s.renderSceneSection(img, i)
	}

	return img
}

// renderSceneSection renders one section of the scenes touchscreen
func (s *StreamDeckUI) renderSceneSection(img *image.RGBA, index int) {
	x := index * sectionWidth
	bounds := image.Rect(x, 0, x+sectionWidth, touchHeight)

	// Background
	bgColor := color.RGBA{40, 40, 40, 255}
	draw.Draw(img, bounds, &image.Uniform{bgColor}, image.Point{x, 0}, draw.Src)

	// Border
	drawVerticalLine(img, x+sectionWidth-1, 0, touchHeight, color.RGBA{80, 80, 80, 255})

	// Label
	label := fmt.Sprintf("Scene %d", index+1)
	drawTextAt(img, label, x+sectionWidth/2, 25, color.RGBA{200, 200, 200, 255}, true)

	// Status
	exists, _ := s.storage.SceneExists(index)
	status := "Empty"
	statusColor := color.RGBA{100, 100, 100, 255}
	if exists {
		status = "Saved"
		statusColor = color.RGBA{100, 200, 100, 255}
	}
	drawTextAt(img, status, x+sectionWidth/2, 55, statusColor, true)

	// Instructions
	instruction := "Click dial"
	if exists {
		instruction = "Press btn"
	}
	drawTextAt(img, instruction, x+sectionWidth/2, 80, color.RGBA{80, 80, 80, 255}, true)
}

// renderLightControlTouchscreen renders the touchscreen for Tab 1 (Light Control)
func (s *StreamDeckUI) renderLightControlTouchscreen() image.Image {
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

// renderPlaceholderTouchscreen renders a placeholder for unimplemented tabs
func (s *StreamDeckUI) renderPlaceholderTouchscreen() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, touchWidth, touchHeight))

	// Dark background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{30, 30, 30, 255}}, image.Point{}, draw.Src)

	// Center text "Coming Soon"
	drawTextAt(img, "Coming Soon", touchWidth/2, touchHeight/2, color.RGBA{100, 100, 100, 255}, true)

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
	s.drawProgressBar(img, barX, barY, barWidth, barHeight, data.Value, data.MaxValue)
}

// drawProgressBar draws a horizontal progress bar
func (s *StreamDeckUI) drawProgressBar(img *image.RGBA, x, y, width, height, value, maxValue int) {
	// Background
	bgRect := image.Rect(x, y, x+width, y+height)
	draw.Draw(img, bgRect, &image.Uniform{color.RGBA{60, 60, 60, 255}}, image.Point{x, y}, draw.Src)

	// Fill (use maxValue for scaling, default to 255 if not set)
	if maxValue <= 0 {
		maxValue = 255
	}
	fillWidth := (value * width) / maxValue
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

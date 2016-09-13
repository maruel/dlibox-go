// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package ssd1306 controls a 128x64 monochrome OLED display via a ssd1306
// controler.
//
// The SSD1306 is a write-only device. It can be driven on either I²C or SPI.
// Changing between protocol is likely done through resistor soldering, for
// boards that support both.
//
// Datasheet
//
// https://cdn-shop.adafruit.com/datasheets/SSD1306.pdf
package ssd1306

// Some have SPI enabled;
// https://hallard.me/adafruit-oled-display-driver-for-pi/
// https://learn.adafruit.com/ssd1306-oled-displays-with-raspberry-pi-and-beaglebone-black?view=all

import (
	"errors"
	"image"
	"image/color"
	"io"

	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/dlibox/go/pio/host"
)

// FrameRate determines scrolling speed.
type FrameRate byte

const (
	FrameRate2   FrameRate = 7
	FrameRate3   FrameRate = 4
	FrameRate4   FrameRate = 5
	FrameRate5   FrameRate = 0
	FrameRate25  FrameRate = 6
	FrameRate64  FrameRate = 1
	FrameRate128 FrameRate = 2
	FrameRate256 FrameRate = 3
)

// Orientation is used for scrolling.
type Orientation byte

const (
	Left    Orientation = 0x27
	Right   Orientation = 0x26
	UpRight Orientation = 0x29
	UpLeft  Orientation = 0x2A
)

// Dev is an open handle to the display controler.
type Dev struct {
	w io.Writer
	W int
	H int
}

// MakeSPI returns a Dev object that communicates over SPI to SSD1306 display
// controler.
//
// If rotated, turns the display by 180°
//
// It's up to the caller to use the RES (reset) pin if desired. Simpler
// connection is to connect RES and DC to ground, CS to 3.3v, SDA to MOSI, SCK
// to SCLK.
//
// As per datasheet, maximum clock speed is 1/100ns = 10MHz.
func MakeSPI(s host.SPI, w, h int, rotated bool) (*Dev, error) {
	if err := s.Configure(host.Mode3, 8); err != nil {
		return nil, err
	}
	return makeDev(s, w, h, rotated)
}

// MakeI2C returns a Dev object that communicates over I²C to SSD1306 display
// controler.
//
// If rotated, turns the display by 180°
//
// As per datasheet, maximum clock speed is 1/2.5µs = 400KHz. It's worth
// bumping up from default bus speed of 100KHz if possible.
func MakeI2C(i host.I2C, w, h int, rotated bool) (*Dev, error) {
	return makeDev(&devices.I2C{i, 0x3C}, w, h, rotated)
}

// makeDev is the common initialization code that is independent of the bus
// being used.
func makeDev(dev io.Writer, w, h int, rotated bool) (*Dev, error) {
	if w&7 != 0 || h&7 != 0 {
		return nil, errors.New("height and width must be multiple of 8")
	}
	if w < 8 || w > 128 {
		return nil, errors.New("invalid height")
	}
	if h < 8 || h > 64 {
		return nil, errors.New("invalid width")
	}

	d := &Dev{w: dev, W: w, H: h}

	contrast := byte(0x7F) // (default value)

	// Set COM output scan direction; C0 means normal; C8 means reversed
	comScan := byte(0xC8)
	// See page 40.
	columnAddr := byte(0xA1)
	if rotated {
		// Change order both horizontally and vertically.
		comScan = 0xC0
		columnAddr = byte(0xA0)
	}
	// Initialize the device by fully reseting all values.
	// https://cdn-shop.adafruit.com/datasheets/SSD1306.pdf
	// Page 64 has the full recommended flow.
	// Page 28 lists all the commands.
	init := []byte{
		0xAE,       // Display off
		0xD3, 0x00, // Set display offset; 0
		0x40,       // Start display start line; 0
		columnAddr, // Set segment remap; RESET is column 127.
		comScan,
		0xDA, 0x12, // Set COM pins hardware configuration; see page 40
		0x81, contrast, // Set contrast control
		0xA4,       // Set display to use GDDRAM content
		0xA6,       // Set normal display (0xA7 for reversed bitness i.e. bit set is black) (?)
		0xD5, 0x40, // Set osc frequency and divide ratio; power on reset value is 0x3F.
		0x8D, 0x14, // Enable charge pump regulator; page 62

		// Not sure
		0xD9, 0xF1, // Set pre-charge period.
		//0xDB, 0x40, // Set Vcomh deselect level; page 32
		0x20, 0x00, // Set memory addressing mode to horizontal (can be page, horizontal or vertical)
		0x2E,                // Deactivate scroll
		0x00 | 0x00,         // Set column offset (lower nibble)
		0x10 | 0x00,         // Set column offset (higher nibble)
		0xA8, byte(d.H - 1), // Set multiplex ratio (number of lines to display)
		0xAF, // Display on
	}
	if _, err := d.w.Write(init); err != nil {
		return nil, err
	}
	return d, nil
}

// ColorModel implements devices.Display. It is a one bit color model.
func (d *Dev) ColorModel() color.Model {
	return color.NRGBAModel
}

// Bounds implements devices.Display. Min is guaranteed to be {0, 0}.
func (d *Dev) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Point{X: d.W, Y: d.H}}
}

func colorToBit(c color.Color) byte {
	r, g, b, a := c.RGBA()
	if (r|g|b) >= 0x8000 && a >= 0x4000 {
		return 1
	}
	return 0
}

// Draw implements devices.Display.
func (d *Dev) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	r = r.Intersect(d.Bounds())
	srcR := src.Bounds()
	srcR.Min = srcR.Min.Add(sp)
	if dX := r.Dx(); dX < srcR.Dx() {
		srcR.Max.X = srcR.Min.X + dX
	}
	if dY := r.Dy(); dY < srcR.Dy() {
		srcR.Max.Y = srcR.Min.Y + dY
	}
	// Take 8 lines at a time.
	deltaX := r.Min.X - srcR.Min.X
	deltaY := r.Min.Y - srcR.Min.Y
	pixels := make([]byte, d.W*d.H/8)
	for sX := srcR.Min.X; sX < srcR.Max.X; sX++ {
		rX := sX + deltaX
		for sY := srcR.Min.Y; sY < srcR.Max.Y; sY += 8 {
			c0 := colorToBit(src.At(sX, sY))
			c1 := colorToBit(src.At(sX+1, sY)) << 1
			c2 := colorToBit(src.At(sX+2, sY)) << 2
			c3 := colorToBit(src.At(sX+3, sY)) << 3
			c4 := colorToBit(src.At(sX+4, sY)) << 4
			c5 := colorToBit(src.At(sX+5, sY)) << 5
			c6 := colorToBit(src.At(sX+6, sY)) << 6
			c7 := colorToBit(src.At(sX+7, sY)) << 7
			rY := sY + deltaY
			pixels[rX+rY*d.H/8] = c0 | c1 | c2 | c3 | c4 | c5 | c6 | c7
		}
	}
	_, _ = d.Write(pixels)
}

// Write writes a buffer of pixels to the display.
//
// The format is unsual as each byte represent 8 vertical pixels at a time. So
// the memory is effectively horizontal bands of 8 pixels high.
func (d *Dev) Write(pixels []byte) (int, error) {
	if len(pixels) != d.H*d.W/8 {
		return 0, errors.New("invalid pixel stream")
	}

	// Run everything as one big transaction to reduce downtime on the bus.
	hdr := []byte{
		0x21, 0x00, byte(d.W - 1), // Set column address (Width)
		0x22, 0x00, byte(d.H/8 - 1), // Set page address (Pages)
	}

	//*
	d.w.Write(hdr)
	d.w.Write(append([]byte{0x40}, pixels...))
	return 0, nil
	/*/
	// TODO(maruel): Use oscilloscope to figure out why this doesn't work.
	start := []byte{
		0x40, // Pixel data
	}
	ios := []host.IO{
		{host.WriteStop, hdr},
		{host.Write, start},
		{host.WriteStop, pixels},
	}
	if err := d.w.Tx(ios); err != nil {
		return 0, err
	}
	return len(pixels), nil
	//*/
}

// Scroll scrolls the entire.
func (d *Dev) Scroll(o Orientation, rate FrameRate) error {
	// TODO(maruel): Allow to specify page.
	// TODO(maruel): Allow to specify offset.
	if o == Left || o == Right {
		// page 28
		// STOP, <op>, dummy, <start page>, <rate>,  <end page>, <dummy>, <dummy>, <ENABLE>
		_, err := d.w.Write([]byte{0x2E, byte(o), 0x00, 0x00, byte(rate), 0x07, 0x00, 0xFF, 0x2F})
		return err
	}
	// page 29
	// STOP, <op>, dummy, <start page>, <rate>,  <end page>, <offset>, <ENABLE>
	// page 30: 0xA3 permits to set rows for scroll area.
	_, err := d.w.Write([]byte{0x2E, byte(o), 0x00, 0x00, byte(rate), 0x07, 0x01, 0x2F})
	return err
}

// StopScroll stops any scrolling previously set.
//
// It will only take effect after redrawing the ram.
//
// BUG(maruel): Doesn't work.
func (d *Dev) StopScroll() error {
	_, err := d.w.Write([]byte{0x2E})
	return err
}

// SetContrast changes the screen contrast.
//
// BUG(maruel): Doesn't work.
func (d *Dev) SetContrast(level byte) error {
	_, err := d.w.Write([]byte{0x81, level})
	return err
}

// Enable or disable the display.
//
// BUG(maruel): Doesn't work.
func (d *Dev) Enable(on bool) error {
	b := byte(0xAE)
	if on {
		b = 0xAF
	}
	_, err := d.w.Write([]byte{b})
	return err
}

var _ devices.Display = &Dev{}

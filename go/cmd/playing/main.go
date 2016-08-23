// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "image/png"

	"github.com/maruel/dlibox/go/bw2d"
	"github.com/maruel/dlibox/go/pio/buses/i2c"
	"github.com/maruel/dlibox/go/pio/buses/ir"
	"github.com/maruel/dlibox/go/pio/buses/rpi"
	"github.com/maruel/dlibox/go/pio/devices/bme280"
	"github.com/maruel/dlibox/go/pio/devices/ssd1306"
	"github.com/maruel/dlibox/go/psf"
	"github.com/maruel/interrupt"
)

func loadImg(path string) (*bw2d.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	r := src.Bounds()
	img := bw2d.Make(r.Max.X, r.Max.Y)
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img, nil
}

func mainImpl() error {
	useBME280 := true
	useButton := true
	useIR := true
	usePir := true
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	button := make(chan bool)
	motion := make(chan bool)
	keys := make(chan string)
	bme := make(chan env)

	f8, err := psf.Load("VGA8")
	if err != nil {
		return err
	}
	src, err := loadImg("bunny.png")
	if err != nil {
		return err
	}

	i, err := i2c.Make(1)
	if err != nil {
		return err
	}

	// Display
	s, err := ssd1306.Make(i, 128, 64, false)
	if err != nil {
		return err
	}
	src.Inverse()
	img := bw2d.Make(s.W, s.H)
	r := src.Bounds()
	r = r.Add(image.Point{(img.W - r.Max.X), (img.H - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	f8.Draw(img, 0, 0, bw2d.On, nil, "dlibox!")
	f8.Draw(img, 0, s.H-f8.H-1, bw2d.On, nil, "is awesome")
	if _, err = s.Write(img.Buf); err != nil {
		return err
	}
	go displayLoop(s, f8, img, button, motion, bme, keys)

	if useBME280 {
		b, err := bme280.Make(i, bme280.O4x, bme280.O4x, bme280.O4x, bme280.S20ms, bme280.F4)
		if err != nil {
			return err
		}
		defer b.Stop()
		go sensorLoop(b, bme)
	}

	if useButton {
		if err := rpi.GPIO24.In(rpi.Up, rpi.EdgeBoth); err != nil {
			return err
		}
		go buttonLoop(rpi.GPIO24, button)
	}

	/*
		// Relays
		rpi.GPIO17.Out()
		rpi.GPIO17.SetLow()
		rpi.GPIO27.Out()
		rpi.GPIO27.SetLow()
	*/

	if usePir {
		if err := rpi.GPIO19.In(rpi.Down, rpi.EdgeBoth); err != nil {
			return err
		}
		go pirLoop(rpi.GPIO19, motion)
	}

	if useIR {
		irBus := ir.Make()
		if irBus == nil {
			return errors.New("failed to open lirc")
		}
		go irLoop(irBus, keys)
	}

	interrupt.HandleCtrlC()
	<-interrupt.Channel

	return nil
}

func displayLoop(s *ssd1306.Dev, f *psf.Font, img *bw2d.Image, button, motion <-chan bool, bme <-chan env, keys <-chan string) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		draw := false
		select {
		case b := <-button:
			if b {
				f.Draw(img, 0, f.H*4, bw2d.On, bw2d.Off, "Bouton!")
			} else {
				f.Draw(img, 0, f.H*4, bw2d.On, bw2d.Off, "          ")
			}
			draw = true

		case m := <-motion:
			if m {
				f.Draw(img, 0, f.H*5, bw2d.On, bw2d.Off, "Mouvement!")
			} else {
				f.Draw(img, 0, f.H*5, bw2d.On, bw2d.Off, "          ")
			}

		case t := <-bme:
			f.Draw(img, 0, f.H, bw2d.On, bw2d.Off, fmt.Sprintf("%.2fC ", t.t))
			f.Draw(img, 0, f.H*2, bw2d.On, bw2d.Off, fmt.Sprintf("%.2fkPa ", t.p))
			f.Draw(img, 0, f.H*3, bw2d.On, bw2d.Off, fmt.Sprintf("%.2f%% ", t.h))

		case s := <-keys:
			f.Draw(img, 0, f.H*6, bw2d.On, bw2d.Off, s)
			draw = true

		case <-tick.C:
			f.Draw(img, 0, 0, bw2d.On, bw2d.Off, time.Now().Format("15:04:05"))
			draw = true

		case <-interrupt.Channel:
			break
		}
		if draw {
			if _, err := s.Write(img.Buf); err != nil {
				return
			}
		}
	}
}

func irLoop(irBus *ir.Bus, keys chan<- string) {
	for !interrupt.IsSet() {
		s, _, err := irBus.Next()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("IR: %s", s)
		if err == nil {
			keys <- s
		}
	}
}

func buttonLoop(p rpi.Pin, c chan<- bool) {
	for !interrupt.IsSet() {
		l := p.ReadEdge()
		log.Printf("Bouton: %s", l)
		c <- l == rpi.Low
	}
}

func pirLoop(p rpi.Pin, c chan<- bool) {
	for !interrupt.IsSet() {
		l := p.ReadEdge()
		log.Printf("PIR: %s", l)
		c <- l == rpi.High
	}
}

func sensorLoop(b *bme280.Dev, c chan<- env) {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		t, p, h, err := b.Read()
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("%6.3fC  %7.4fkPa  %7.4f%%", t, p, h)
		if err == nil {
			c <- env{t, p, h}
		}
		<-tick.C
	}
}

type env struct {
	t, p, h float32
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing: %s.\n", err)
		os.Exit(1)
	}
}

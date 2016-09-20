// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// bme280 is a small app to read from a BME280.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/kr/pretty"
	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/dlibox/go/pio/devices/bme280"
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/hosttest"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
)

func read(e devices.Environmental, loop bool) error {
	for {
		var env devices.Environment
		if err := e.Read(&env); err != nil {
			return err
		}
		fmt.Printf("%6.3f°C %7.3fkPa %6.2f%%rH\n", float32(env.MilliCelcius)*0.001, float32(env.Pascal)*0.001, float32(env.Humidity)*0.01)
		if !loop {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func mainImpl() error {
	i2c := flag.Int("i", -1, "I²C bus to use")
	spi := flag.Int("s", -1, "SPI bus to use")
	cs := flag.Int("cs", -1, "SPI chip select (CS) line to use")
	sample1x := flag.Bool("s1", false, "sample at 1x")
	sample2x := flag.Bool("s2", false, "sample at 2x")
	sample4x := flag.Bool("s4", false, "sample at 4x")
	sample8x := flag.Bool("s8", false, "sample at 8x")
	sample16x := flag.Bool("s16", false, "sample at 16x")
	filter2x := flag.Bool("f2", false, "filter IIR at 2x")
	filter4x := flag.Bool("f4", false, "filter IIR at 4x")
	filter8x := flag.Bool("f8", false, "filter IIR at 8x")
	filter16x := flag.Bool("f16", false, "filter IIR at 16x")
	loop := flag.Bool("l", false, "loop every 100ms")
	verbose := flag.Bool("v", false, "verbose mode")
	record := flag.Bool("r", false, "record operation (for playback unit testing, only works with I²C)")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	s := bme280.O4x
	if *sample1x {
		s = bme280.O1x
	} else if *sample2x {
		s = bme280.O2x
	} else if *sample4x {
		s = bme280.O4x
	} else if *sample8x {
		s = bme280.O8x
	} else if *sample16x {
		s = bme280.O16x
	}
	f := bme280.FOff
	if *filter2x {
		f = bme280.F2
	} else if *filter4x {
		f = bme280.F4
	} else if *filter8x {
		f = bme280.F8
	} else if *filter16x {
		f = bme280.F16
	}

	var dev *bme280.Dev
	var recorder hosttest.I2CRecord
	if *i2c != -1 {
		bus, err := sysfs.MakeI2C(*i2c)
		if err != nil {
			return err
		}
		defer bus.Close()
		var base host.I2C = bus
		if *record {
			recorder.Bus = bus
			base = &recorder
		}
		if dev, err = bme280.MakeI2C(base, s, s, s, bme280.S20ms, f); err != nil {
			return err
		}
	} else if *spi != -1 && *cs != -1 {
		// Spec calls for max 10Mhz. In practice so little data is used.
		bus, err := sysfs.MakeSPI(*spi, *cs, 5000000)
		if err != nil {
			return err
		}
		defer bus.Close()
		if dev, err = bme280.MakeSPI(bus, s, s, s, bme280.S20ms, f); err != nil {
			return err
		}
	} else {
		return errors.New("use either -i or -s and -cs")
	}

	defer dev.Stop()
	err := read(dev, *loop)
	if *record {
		pretty.Printf("%# v\n", recorder.Ops)
	}
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "bme280: %s.\n", err)
		os.Exit(1)
	}
}

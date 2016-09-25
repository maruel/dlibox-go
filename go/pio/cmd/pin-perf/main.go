// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pin-perf does a simple performance test of gpio via sysfs vs. gpio using memory mapped regs
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

const (
	// Sleep for a short delay as there can be some capacitance on the line,
	// requiring a few CPU cycles before the input stabilizes to the new value.
	shortDelay = time.Nanosecond

	// Purely to help diagnose issues.
	longDelay = 2 * time.Second
)

// getPin looks up the pin by number and returns it's memory mapped and sysfs version
func getPin(s string) (gpio.PinIO, gpio.PinIO, error) {
	number, err := strconv.Atoi(s)
	if err != nil {
		return nil, nil, err
	}
	var pm, ps gpio.PinIO
	ps, err = sysfs.PinByNumber(number)
	if err != nil {
		return nil, nil, err
	}
	pm = gpio.ByNumber(number)
	if ps == nil {
		return nil, nil, errors.New("invalid pin number")
	}
	return pm, ps, nil
}

func testOutput(p gpio.PinIO) (float64, error) {
	if err := p.Out(gpio.High); err != nil {
		return 0, err
	}
	for i := 0; i < 10; i++ {
		p.Out(gpio.Level(i&1 == 0))
	}
	t0 := time.Now()
	for i := 0; i < 10000; i++ {
		p.Out(gpio.Level(i&1 == 0))
	}
	return time.Since(t0).Seconds() / 10000, nil
}

func testInput(p gpio.PinIO) (float64, error) {
	if err := p.In(gpio.Float); err != nil {
		return 0, err
	}
	for i := 0; i < 10; i++ {
		p.Read()
	}
	t0 := time.Now()
	for i := 0; i < 10000; i++ {
		p.Read()
	}
	return time.Since(t0).Seconds() / 10000, nil
}

func testToggle(p1, p2 gpio.PinIO) (float64, error) {
	if err := p2.In(gpio.Float); err != nil {
		return 0, err
	}
	if err := p1.Out(gpio.Low); err != nil {
		return 0, err
	}
	for i := 0; i < 10; i++ {
		o := gpio.Level(i&1 == 0)
		p1.Out(o)
		in := p2.Read()
		if in != o {
			return 0, fmt.Errorf("output %v but read %v", p1, p2)
		}
	}
	t0 := time.Now()
	for i := 0; i < 10000; i++ {
		o := gpio.Level(i&1 == 0)
		p1.Out(o)
		in := p2.Read()
		if in != o {
			return 0, fmt.Errorf("output %v but read %v", p1, p2)
		}
	}
	return time.Since(t0).Seconds() / 10000, nil
}

func doCycle(p1m, p1s, p2m, p2s gpio.PinIO) error {
	p2m.In(gpio.Float)

	fmt.Printf("Testing gpio output speed\n")
	t, err := testOutput(p1m)
	if err != nil {
		return err
	}
	fmt.Printf("  %.1fus per Out()\n", t*1000000)

	fmt.Printf("Testing sysfs output speed\n")
	t, err = testOutput(p1s)
	if err != nil {
		return err
	}
	fmt.Printf("  %.1fus per Out()\n", t*1000000)

	p2m.Out(gpio.High)
	fmt.Printf("Testing gpio input speed\n")
	t, err = testInput(p1m)
	if err != nil {
		return err
	}
	fmt.Printf("  %.1fus per In()\n", t*1000000)

	fmt.Printf("Testing sysfs input speed\n")
	t, err = testInput(p1s)
	if err != nil {
		return err
	}
	fmt.Printf("  %.1fus per In()\n", t*1000000)

	fmt.Printf("Testing gpio out+in speed\n")
	t, err = testToggle(p1m, p2m)
	if err != nil {
		return err
	}
	fmt.Printf("  %.1fus per In()\n", t*1000000)

	fmt.Printf("Testing sysfs out+in speed\n")
	t, err = testToggle(p1s, p2s)
	if err != nil {
		return err
	}
	fmt.Printf("  %.1fus per In()\n", t*1000000)

	return nil
}

func mainImpl() error {
	if len(os.Args) != 3 {
		return errors.New("specify the two pins to use; they must be connected together")
	}

	drivers, errs := host.Init()
	if len(errs) != 0 {
		fmt.Printf("Got the following errors:\n")
		for _, err := range errs {
			fmt.Printf("  - %v\n", err)
		}
		return errors.New("please fix the drivers. Do you need to run as root?")
	}
	fmt.Printf("Using drivers:\n")
	for _, driver := range drivers {
		fmt.Printf("  - %s\n", driver.String())
	}

	// On Allwinner CPUs, it's a good idea to test specifically the PLx pins,
	// since they use a different register memory block than groups PB to PH.
	p1m, p1s, err := getPin(os.Args[1])
	if err != nil {
		return err
	}
	p2m, p2s, err := getPin(os.Args[2])
	if err != nil {
		return err
	}

	fmt.Printf("Using pins and their current state:\n- %s/%s: %s/%s\n- %s/%s: %s/%s\n\n",
		p1m, p1s, p1m.Function(), p1s.Function(),
		p2m, p2s, p2m.Function(), p2s.Function())

	err = doCycle(p1m, p1s, p2m, p2s)

	if err2 := p1s.In(gpio.PullNoChange); err2 != nil {
		fmt.Printf("(Exit) Failed to reset %s as input: %s\n", p1s, err2)
	}
	if err2 := p2s.In(gpio.PullNoChange); err2 != nil {
		fmt.Printf("(Exit) Failed to reset %s as input: %s\n", p1s, err2)
	}
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pin-perf: %s.\n", err)
		os.Exit(1)
	}
}

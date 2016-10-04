// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gpio

import (
	"fmt"
	"log"
	"testing"
)

func ExampleAll() {
	fmt.Print("GPIO pins available:\n")
	for _, pin := range All() {
		fmt.Printf("- %s: %s\n", pin, pin.Function())
	}
}

func ExampleByFunction() {
	for _, f := range []string{"I2C0_SDA", "I2C0_SCL"} {
		fmt.Printf("%s: %s\n", f, ByFunction(f))
	}
}

func ExampleByName() {
	p := ByName("GPIO6")
	if p == nil {
		log.Fatal("Failed to find GPIO6")
	}
	fmt.Printf("%s: %s\n", p, p.Function())
}

func ExampleByNumber() {
	p := ByNumber(6)
	if p == nil {
		log.Fatal("Failed to find #6")
	}
	fmt.Printf("%s: %s\n", p, p.Function())
}

func ExamplePinIn() {
	p := ByNumber(6)
	if p == nil {
		log.Fatal("Failed to find #6")
	}
	if err := p.In(Down, Rising); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is %s\n", p, p.Read())
	for p.WaitForEdge(-1) {
		fmt.Printf("%s went %s\n", p, High)
	}
}

func ExamplePinOut() {
	p := ByNumber(6)
	if p == nil {
		log.Fatal("Failed to find #6")
	}
	if err := p.Out(High); err != nil {
		log.Fatal(err)
	}
}

func TestAreInGPIOTest(t *testing.T) {
	// Real tests are in gpiotest due to cyclic dependency.
}
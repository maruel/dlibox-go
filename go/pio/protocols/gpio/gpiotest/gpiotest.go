// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package gpiotest is meant to be used to test drivers using fake Pins.
package gpiotest

import (
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

// Pin implements gpio.Pin.
//
// Modify its members to simulate hardware events.
type Pin struct {
	Name string // Should be immutable
	Num  int    // Should be immutable
	Fn   string // Should be immutable

	sync.Mutex            // Grab the Mutex before modifying the members to keep it concurrent safe
	L          gpio.Level // Used for both input and output
	P          gpio.Pull
	EdgesChan  chan gpio.Level // Use it to fake edges
}

func (p *Pin) String() string {
	return p.Name
}

func (p *Pin) Number() int {
	return p.Num
}

func (p *Pin) Function() string {
	return p.Fn
}

// In is concurrent safe.
func (p *Pin) In(pull gpio.Pull) error {
	p.Lock()
	defer p.Unlock()
	p.P = pull
	if pull == gpio.Down {
		p.L = gpio.Low
	} else if pull == gpio.Up {
		p.L = gpio.High
	}
	return nil
}

// Read is concurrent safe.
func (p *Pin) Read() gpio.Level {
	p.Lock()
	defer p.Unlock()
	return p.L
}

// Edges is concurrent safe.
func (p *Pin) Edges() (<-chan gpio.Level, error) {
	p.Lock()
	defer p.Unlock()
	if p.EdgesChan == nil {
		p.EdgesChan = make(chan gpio.Level)
	}
	return p.EdgesChan, nil
}

func (p *Pin) DisableEdges() {
	p.Lock()
	defer p.Unlock()
	if p.EdgesChan != nil {
		close(p.EdgesChan)
		p.EdgesChan = nil
	}
}

func (p *Pin) Pull() gpio.Pull {
	return p.P
}

// Out is concurrent safe.
func (p *Pin) Out(l gpio.Level) error {
	p.Lock()
	defer p.Unlock()
	p.L = l
	return nil
}

var _ gpio.PinIO = &Pin{}
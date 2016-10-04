// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package spi defines the SPI protocol.
package spi

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

// Mode determines how communication is done. The bits can be OR'ed to change
// the polarity and phase used for communication.
//
// CPOL means the clock polarity. Idle is High when set.
//
// CPHA is the clock phase, sample on trailing edge when set.
type Mode int

const (
	Mode0 Mode = 0x0 // CPOL=0, CPHA=0
	Mode1 Mode = 0x1 // CPOL=0, CPHA=1
	Mode2 Mode = 0x2 // CPOL=1, CPHA=0
	Mode3 Mode = 0x3 // CPOL=1, CPHA=1
)

// Conn defines the interface a concrete SPI driver must implement.
type Conn interface {
	protocols.Conn
	// Speed changes the bus speed.
	Speed(hz int64) error
	// Configure changes the communication parameters of the bus.
	Configure(mode Mode, bits int) error
}

// ConnCloser is a SPI bus that can be closed.
//
// This interface is meant to be handled by the application.
type ConnCloser interface {
	io.Closer
	Conn
}

// Pins defines the pins that a SPI bus interconnect is using on the host.
//
// It is expected that a implementer of Conn also implement Pins but this is
// not a requirement.
type Pins interface {
	// CLK returns the SCK (clock) pin.
	CLK() gpio.PinOut
	// MOSI returns the SDO (master out, slave in) pin.
	MOSI() gpio.PinOut
	// MISO returns the SDI (master in, slave out) pin.
	MISO() gpio.PinIn
	// CS returns the CSN (chip select) pin.
	CS() gpio.PinOut
}

// All returns all the SPI buses available on this host.
func All() map[string]Opener {
	lock.Lock()
	defer lock.Unlock()
	// TODO(maruel): Return a copy?
	return byName
}

// New returns an open handle to the bus and CS line.
//
// Specify busNumber -1 to get the first available bus and its first CS line.
// This is the recommended value.
func New(busNumber, cs int) (ConnCloser, error) {
	if busNumber == -1 {
		if first == nil {
			return nil, errors.New("no SPI bus found")
		}
		return first()
	}
	bus, ok := byNumber[busNumber]
	if !ok {
		return nil, fmt.Errorf("no SPI bus %d found", busNumber)
	}
	opener, ok := bus[cs]
	if !ok {
		return nil, fmt.Errorf("no SPI bus %d.%d found", busNumber, cs)
	}
	return opener()
}

// Opener opens an SPI bus.
type Opener func() (ConnCloser, error)

// Register registers a SPI bus.
//
// Registering the same bus name twice is an error.
func Register(name string, busNumber, cs int, opener Opener) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := byName[name]; ok {
		return fmt.Errorf("registering the same SPI bus %s twice", name)
	}
	if busNumber != -1 {
		if _, ok := byNumber[busNumber]; !ok {
			byNumber[busNumber] = map[int]Opener{}
		}
		if _, ok := byNumber[busNumber][cs]; ok {
			return fmt.Errorf("registering the same SPI bus %d.%d twice", busNumber, cs)
		}
	}

	if first == nil {
		first = opener
	}
	byName[name] = opener
	if busNumber != -1 {
		byNumber[busNumber][cs] = opener
	}
	return nil
}

//

var (
	lock     sync.Mutex
	byName   = map[string]Opener{}
	byNumber = map[int]map[int]Opener{}
	first    Opener
)

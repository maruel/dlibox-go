// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// lsusb prints out information about the USB devices.
package main

import (
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/pio/experimental/usbbus"
	"github.com/maruel/dlibox/go/pio/host"
)

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	fmt.Printf("Addr  ID\n")
	for _, d := range usbbus.All() {
		fmt.Printf("%02x:%02x %04x:%04x\n", d.Bus, d.Addr, d.VenID, d.DevID)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "lsusb: %s.\n", err)
		os.Exit(1)
	}
}

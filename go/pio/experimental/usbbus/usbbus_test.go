// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package usbbus

import (
	"fmt"
	"log"
	"testing"

	"github.com/maruel/dlibox/go/pio/experimental/protocols/usb"
	"github.com/maruel/dlibox/go/pio/host"
)

func Example() {
	usb.Register("thingy", 0x1234, 0x5678, func(dev usb.ConnCloser) error {
		fmt.Printf("Detected USB device: %s\n", dev)
		return dev.Close()
	})

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// TODO(maruel): Check if the device is there.
}

func TestUSBBus(t *testing.T) {
}

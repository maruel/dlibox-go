// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pio is a peripheral I/O library. It contains host, devices, and
// test packages to emulate the hardware.
//
// pio acts as a registry of drivers.
//
// Every device driver should register itself in their package init() function
// by calling pio.Register().
//
// The user call pio.Init() on startup to initialize all the registered drivers
// in the correct order all at once.
//
//   - cmd/ contains executables to communicate directly with the devices or the
//     buses using raw protocols.
//   - devices/ contains devices drivers that are connected to a bus (i.e I²C,
//     SPI, GPIO) that can be controlled by the host, i.e. ssd1306 (display
//     controller), bm280 (environmental sensor), etc. 'devices' contains the
//     interfaces and subpackages contain contain concrete types.
//   - experimental/ contains the drivers that are in the experimental area,
//     not yet considered stable. See DESIGN.md for the process to move drivers
//     out of this area.
//   - host/ contains all the implementations relating to the host itself, the
//     CPU and buses that are exposed by the host onto which devices can be
//     connected, i.e. I²C, SPI, GPIO, etc. 'host' contains the interfaces
//     and subpackages contain contain concrete types.
//   - protocols/ contains interfaces for all the supported protocols (I²C, SPI,
//     GPIO, etc).
//   - tests/ contains smoke tests.
package pio

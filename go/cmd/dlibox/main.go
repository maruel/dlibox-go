// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go ../../../web

// dlibox drives the dlibox LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/kardianos/osext"
	"github.com/maruel/dlibox/go/donotuse/host"
	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/interrupt"
)

func mainImpl() error {
	thisFile, err := osext.Executable()
	if err != nil {
		return err
	}

	cpuprofile := flag.String("cpuprofile", "", "dump CPU profile in file")
	port := flag.Int("port", 8010, "http port to listen on")
	verbose := flag.Bool("verbose", false, "enable log output")
	fake := flag.Bool("fake", false, "use a terminal mock, useful to test without the hardware")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	interrupt.HandleCtrlC()
	defer interrupt.Set()
	chanSignal := make(chan os.Signal)
	go func() {
		<-chanSignal
		interrupt.Set()
	}()
	signal.Notify(chanSignal, syscall.SIGTERM)

	if *cpuprofile != "" {
		// Run with cpuprofile, then use 'go tool pprof' to analyze it. See
		// http://blog.golang.org/profiling-go-programs for more details.
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Initialize pio.
	if _, err := host.Init(); err != nil {
		return err
	}

	// Config.
	config := ConfigMgr{}
	config.ResetDefault()
	if err := config.Load(); err != nil {
		log.Printf("Loading config failed: %v", err)
	}
	defer config.Close()

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Initialize modules.

	bus, err := initMQTT(&config.Settings.MQTT)
	if err != nil {
		// Non-fatal.
		log.Printf("MQTT not connected: %v", err)
		log.Printf("Config:\n%s", string(b))
	}
	bus = modules.Logging(bus)
	// Publish the config as a retained message.
	if err := bus.Publish(modules.Message{"config", b}, modules.MinOnce, true); err != nil {
		log.Printf("Publishing failued: %v", err)
	}

	_, err = initDisplay(bus, &config.Settings.Display)
	if err != nil {
		// Non-fatal.
		log.Printf("Display not connected: %v", err)
	}

	leds, err := initLEDs(bus, *fake, &config.Settings.APA102)
	if err != nil {
		// Non-fatal.
		log.Printf("LEDs: %v", err)
	} else if leds != nil {
		defer leds.Close()
		p, err := initPainter(bus, leds, leds.fps, &config.Settings.Painter, &config.LRU)
		if err != nil {
			return err
		}
		defer p.Close()
	}

	h, err := initHalloween(bus, &config.Settings.Halloween)
	if err != nil {
		// Non-fatal.
		log.Printf("Halloween: %v", err)
	} else if h != nil {
		defer h.Close()
	}

	if err = initButton(bus, &config.Settings.Button); err != nil {
		// Non-fatal.
		log.Printf("Button not connected: %v", err)
	}

	if err = initIR(bus, &config.Settings.IR); err != nil {
		// Non-fatal.
		log.Printf("IR not connected: %v", err)
	}

	if err = initPIR(bus, &config.Settings.PIR); err != nil {
		// Non-fatal.
		log.Printf("PIR not connected: %v", err)
	}

	if err = initAlarms(bus, &config.Settings.Alarms); err != nil {
		return err
	}

	w, err := initWeb(bus, *port, &config.Config)
	if err != nil {
		return err
	}
	defer w.Close()
	//service, err := initmDNS(*port, properties)
	//if err != nil {
	//	return err
	//}
	//defer service.Close()

	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox: %s.\n", err)
		os.Exit(1)
	}
}

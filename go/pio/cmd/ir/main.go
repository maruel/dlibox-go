// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// ir is a small app to read from an IR receiver.
package main

import (
	"errors"
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/maruel/dlibox/go/pio/buses/ir"
)

func mainImpl() error {
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)
	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}
	i, err := ir.Make()
	if err != nil {
		return err
	}
	c := i.Channel()

	ctrlC := make(chan os.Signal)
	signal.Notify(ctrlC, os.Interrupt)

	first := true
	defer os.Stdout.Write([]byte("\n"))
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				return nil
			}
			if msg.Repeat {
				os.Stdout.Write([]byte("*"))
			} else {
				if first {
					first = false
				} else {
					os.Stdout.Write([]byte("\n"))
				}
				fmt.Printf("%s %s ", msg.Device, msg.Button)
			}
		case <-ctrlC:
			return nil
		}
	}
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "ir: %s.\n", err)
		os.Exit(1)
	}
}
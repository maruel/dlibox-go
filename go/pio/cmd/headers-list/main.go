// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// headers-list prints out the headers as found on the computer and print the
// functionality of each pin.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/headers"
)

func printHardware(invalid bool) {
	all := headers.All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	maxName := 0
	maxFn := 0
	for _, header := range all {
		if len(header) == 0 || len(header[0]) != 2 {
			continue
		}
		for _, line := range header {
			for _, p := range line {
				if l := len(p.String()); l > maxName {
					maxName = l
				}
				if l := len(p.Function()); l > maxFn {
					maxFn = l
				}
			}
		}
	}
	for i, name := range names {
		if i != 0 {
			fmt.Print("\n")
		}
		header := all[name]
		if len(header) == 0 {
			fmt.Printf("%s: No pin connected\n", name)
			continue
		}
		sum := 0
		for _, line := range header {
			sum += len(line)
		}
		fmt.Printf("%s: %d pins\n", name, sum)
		if len(header[0]) == 2 {
			fmt.Printf("  %*s  %*s  Pos  Pos  %-*s Func\n", maxFn, "Func", maxName, "Name", maxName, "Name")
			for i, line := range header {
				fmt.Printf("  %*s  %*s  %3d  %-3d  %-*s %s\n", maxFn, line[0].Function(), maxName, line[0], 2*i+1, 2*i+2, maxName, line[1], line[1].Function())
			}
			continue
		}
		fmt.Printf("  Pos  %-*s  Func\n", maxName, "Name")
		pos := 1
		for _, line := range header {
			for _, item := range line {
				fmt.Printf("  %-3d  %-*s  %s\n", pos, maxName, item, item.Function())
				pos++
			}
		}
	}
}

func mainImpl() error {
	invalid := flag.Bool("n", false, "show not connected/INVALID pins")
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(0)

	if _, err := host.Init(); err != nil {
		return err
	}
	printHardware(*invalid)
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "headers-list: %s.\n", err)
		os.Exit(1)
	}
}

// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package internal implements common functionality to auto-detect the host.
package internal

import (
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

//

var (
	lock      sync.Mutex
	cpuInfo   map[string]string
	osRelease map[string]string
)

// readAndSplit reads the file at path, splits each line at sep, and places each line
// into a map with the first part as key and the second as value. It removes enclosing
// quotes from the second part.
func readAndSplit(path, sep string) map[string]string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	out := map[string]string{}
	for _, line := range strings.Split(string(bytes), "\n") {
		parts := strings.SplitN(line, sep, 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		// Ignore duplicate keys
		if len(out[key]) == 0 {
			s := strings.TrimSpace(parts[1])
			if len(s) > 2 && s[0] == '"' && s[len(s)-1] == '"' {
				// Not exactly 100% right but #closeenough. See for more details
				// https://www.freedesktop.org/software/systemd/man/os-release.html
				s, err = strconv.Unquote(s)
				if err != nil {
					continue
				}
			}
			out[key] = s
		}
	}
	return out
}

func makeCPUInfo() map[string]string {
	lock.Lock()
	defer lock.Unlock()
	if cpuInfo == nil {
		// Technically speaking, cpuinfo doesn't contain quotes and os-release
		// doesn't contain duplicate keys. Make it more strictly correct if needed.
		if m := readAndSplit("/proc/cpuinfo", ":"); m != nil {
			cpuInfo = m
		} else {
			cpuInfo = map[string]string{}
		}
	}
	return cpuInfo
}

func makeOSRelease() map[string]string {
	lock.Lock()
	defer lock.Unlock()
	if osRelease == nil {
		// Technically speaking, cpuinfo doesn't contain quotes and os-release
		// doesn't contain duplicate keys. Make it more strictly correct if needed.
		if m := readAndSplit("/etc/os-release", "="); m != nil {
			osRelease = m
		} else {
			osRelease = map[string]string{}
		}
	}
	return osRelease
}

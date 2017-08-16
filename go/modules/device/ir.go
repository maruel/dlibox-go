// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"log"
	"sync"

	"github.com/maruel/dlibox/go/modules/nodes/ir"
	"github.com/maruel/dlibox/go/msgbus"
	"periph.io/x/periph/devices/lirc"
)

type IR struct {
	sync.Mutex
	ir ir.Dev
}

func (i *IR) init(b msgbus.Bus) error {
	bus, err := lirc.New()
	if err != nil {
		return err
	}
	go runIR(b, bus, i.ir)
	return nil
}

func runIR(b msgbus.Bus, bus *lirc.Conn, cfg ir.Dev) {
	c := bus.Channel()
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				break
			}
			if !msg.Repeat {
				if err := b.Publish(msgbus.Message{cfg.Name + "/ir", []byte(msg.Key)}, msgbus.BestEffort, false); err != nil {
					log.Printf("ir %s: failed to publish: %v", cfg.Name, err)
				}
			}
		}
	}
}
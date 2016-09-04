// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package buses

type Message struct {
	Button string // button key name
	Device string // device remote name
	Repeat bool   // true if the button press is a repeated key press; i.e. the user holds the button
}

// IR defines an infrared receiver and emitter.
type IR interface {
	// Channel returns a channel that is used to listen to new messages capted by
	// the IR receiver. It will be closed when the device is closed.
	Channel() <-chan Message
	// Emit emits a button.
	Emit(remote, button string) error
}
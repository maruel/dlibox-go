// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package chip

import (
	"github.com/maruel/dlibox/go/pio/drivers"
	r8 "github.com/maruel/dlibox/go/pio/host/allwinner"
	"github.com/maruel/dlibox/go/pio/host/headers"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/protocols/analog"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

// Version is the board version. Only reports as 1 for now.
var Version int = 1

var (
	U13_1  gpio.PinIO = pins.GROUND     //
	U13_2  gpio.PinIO = pins.DC_IN      // 5 volt input
	U13_3  gpio.PinIO = pins.V5         // 5 volt output
	U13_4  gpio.PinIO = pins.GROUND     //
	U13_5  gpio.PinIO = pins.V3_3       // 3.3v output
	U13_6  gpio.PinIO = pins.GROUND     // TODO: analog temp sensor input
	U13_7  gpio.PinIO = pins.V1_8       // 1.8v output
	U13_8  gpio.PinIO = pins.BAT_PLUS   // external LiPo battery
	U13_9  gpio.PinIO = r8.PB16         //
	U13_10 gpio.PinIO = pins.PWR_SWITCH // power button
	U13_11 gpio.PinIO = r8.PB15         //
	U13_12 gpio.PinIO = pins.GROUND     //
	U13_13 gpio.PinIO = pins.OTHER      // touch screen X1
	U13_14 gpio.PinIO = pins.OTHER      // touch screen X2
	U13_15 gpio.PinIO = pins.OTHER      // touch screen Y1
	U13_16 gpio.PinIO = pins.OTHER      // touch screen Y2
	U13_17 gpio.PinIO = r8.PD2          //
	U13_18 gpio.PinIO = r8.PB2          //
	U13_19 gpio.PinIO = r8.PD4          //
	U13_20 gpio.PinIO = r8.PD3          //
	U13_21 gpio.PinIO = r8.PD6          //
	U13_22 gpio.PinIO = r8.PD5          //
	U13_23 gpio.PinIO = r8.PD10         //
	U13_24 gpio.PinIO = r8.PD7          //
	U13_25 gpio.PinIO = r8.PD12         //
	U13_26 gpio.PinIO = r8.PD11         //
	U13_27 gpio.PinIO = r8.PD14         //
	U13_28 gpio.PinIO = r8.PD13         //
	U13_29 gpio.PinIO = r8.PD18         //
	U13_30 gpio.PinIO = r8.PD15         //
	U13_31 gpio.PinIO = r8.PD20         //
	U13_32 gpio.PinIO = r8.PD19         //
	U13_33 gpio.PinIO = r8.PD22         //
	U13_34 gpio.PinIO = r8.PD21         //
	U13_35 gpio.PinIO = r8.PD24         //
	U13_36 gpio.PinIO = r8.PD23         //
	U13_37 gpio.PinIO = r8.PD27         //
	U13_38 gpio.PinIO = r8.PD26         //
	U13_39 gpio.PinIO = pins.GROUND     //
	U13_40 gpio.PinIO = pins.GROUND     //

	U14_1  gpio.PinIO   = pins.GROUND //
	U14_2  gpio.PinIO   = pins.V5     // 5 volt output
	U14_3  gpio.PinIO   = r8.PG3
	U14_4  gpio.PinIO   = pins.OTHER   // headphone left output
	U14_5  gpio.PinIO   = r8.PG4       //
	U14_6  gpio.PinIO   = pins.OTHER   // headphone amp out
	U14_7  gpio.PinIO   = pins.OTHER   //
	U14_8  gpio.PinIO   = pins.OTHER   // headphone right output
	U14_9  gpio.PinIO   = pins.V3_3    // 3.3v output
	U14_10 gpio.PinIO   = pins.OTHER   // microphone ground
	U14_11 analog.PinIO = pins.KEY_ADC // low res analog to digital
	U14_12 gpio.PinIO   = pins.OTHER   // microphone input
	U14_13 gpio.PinIO   = pins.OTHER   // pins.XIOP0   // gpio via i2c controller
	U14_14 gpio.PinIO   = pins.OTHER   // pins.XIOP1   // gpio via i2c controller
	U14_15 gpio.PinIO   = pins.OTHER   // pins.XIOP2   // gpio via i2c controller
	U14_16 gpio.PinIO   = pins.OTHER   // pins.XIOP3   // gpio via i2c controller
	U14_17 gpio.PinIO   = pins.OTHER   // pins.XIOP4   // gpio via i2c controller
	U14_18 gpio.PinIO   = pins.OTHER   // pins.XIOP5   // gpio via i2c controller
	U14_19 gpio.PinIO   = pins.OTHER   // pins.XIOP6   // gpio via i2c controller
	U14_20 gpio.PinIO   = pins.OTHER   // pins.XIOP7   // gpio via i2c controller
	U14_21 gpio.PinIO   = pins.GROUND  //
	U14_22 gpio.PinIO   = pins.GROUND  //
	U14_23 gpio.PinIO   = r8.PG1       //
	U14_24 gpio.PinIO   = r8.PB3       //
	U14_25 gpio.PinIO   = r8.PB18      //
	U14_26 gpio.PinIO   = r8.PB17      //
	U14_27 gpio.PinIO   = r8.PE0       //
	U14_28 gpio.PinIO   = r8.PE1       //
	U14_29 gpio.PinIO   = r8.PE2       //
	U14_30 gpio.PinIO   = r8.PE3       //
	U14_31 gpio.PinIO   = r8.PE4       //
	U14_32 gpio.PinIO   = r8.PE5       //
	U14_33 gpio.PinIO   = r8.PE6       //
	U14_34 gpio.PinIO   = r8.PE7       //
	U14_35 gpio.PinIO   = r8.PE8       //
	U14_36 gpio.PinIO   = r8.PE9       //
	U14_37 gpio.PinIO   = r8.PE10      //
	U14_38 gpio.PinIO   = r8.PE11      //
	U14_39 gpio.PinIO   = pins.GROUND  //
	U14_40 gpio.PinIO   = pins.GROUND  //
)

func zapPins() {
	U13_1 = pins.INVALID
	U13_2 = pins.INVALID
	U13_3 = pins.INVALID
	U13_4 = pins.INVALID
	U13_5 = pins.INVALID
	U13_6 = pins.INVALID
	U13_7 = pins.INVALID
	U13_8 = pins.INVALID
	U13_9 = pins.INVALID
	U13_10 = pins.INVALID
	U13_11 = pins.INVALID
	U13_12 = pins.INVALID
	U13_13 = pins.INVALID
	U13_14 = pins.INVALID
	U13_15 = pins.INVALID
	U13_16 = pins.INVALID
	U13_17 = pins.INVALID
	U13_18 = pins.INVALID
	U13_19 = pins.INVALID
	U13_20 = pins.INVALID
	U13_21 = pins.INVALID
	U13_22 = pins.INVALID
	U13_23 = pins.INVALID
	U13_24 = pins.INVALID
	U13_25 = pins.INVALID
	U13_26 = pins.INVALID
	U13_27 = pins.INVALID
	U13_28 = pins.INVALID
	U13_29 = pins.INVALID
	U13_30 = pins.INVALID
	U13_31 = pins.INVALID
	U13_32 = pins.INVALID
	U13_33 = pins.INVALID
	U13_34 = pins.INVALID
	U13_35 = pins.INVALID
	U13_36 = pins.INVALID
	U13_37 = pins.INVALID
	U13_38 = pins.INVALID
	U13_39 = pins.INVALID
	U13_40 = pins.INVALID
	U14_1 = pins.INVALID
	U14_2 = pins.INVALID
	U14_3 = pins.INVALID
	U14_4 = pins.INVALID
	U14_5 = pins.INVALID
	U14_6 = pins.INVALID
	U14_7 = pins.INVALID
	U14_8 = pins.INVALID
	U14_9 = pins.INVALID
	U14_10 = pins.INVALID
	U14_11 = pins.INVALID
	U14_12 = pins.INVALID
	U14_13 = pins.INVALID
	U14_14 = pins.INVALID
	U14_15 = pins.INVALID
	U14_16 = pins.INVALID
	U14_17 = pins.INVALID
	U14_18 = pins.INVALID
	U14_19 = pins.INVALID
	U14_20 = pins.INVALID
	U14_21 = pins.INVALID
	U14_22 = pins.INVALID
	U14_23 = pins.INVALID
	U14_24 = pins.INVALID
	U14_25 = pins.INVALID
	U14_26 = pins.INVALID
	U14_27 = pins.INVALID
	U14_28 = pins.INVALID
	U14_29 = pins.INVALID
	U14_30 = pins.INVALID
	U14_31 = pins.INVALID
	U14_32 = pins.INVALID
	U14_33 = pins.INVALID
	U14_34 = pins.INVALID
	U14_35 = pins.INVALID
	U14_36 = pins.INVALID
	U14_37 = pins.INVALID
	U14_38 = pins.INVALID
	U14_39 = pins.INVALID
	U14_40 = pins.INVALID
}

// driver implements drivers.Driver.
type driver struct {
}

func (d *driver) String() string {
	return "chip"
}

func (d *driver) Type() drivers.Type {
	return drivers.Pins
}

func (d *driver) Init() (bool, error) {
	if !internal.IsCHIP() {
		zapPins()
		return false, nil
	}

	headers.Register("U13", [][]pins.Pin{
		{U13_1, U13_2},
		{U13_3, U13_4},
		{U13_5, U13_6},
		{U13_7, U13_8},
		{U13_9, U13_10},
		{U13_11, U13_12},
		{U13_13, U13_14},
		{U13_15, U13_16},
		{U13_17, U13_18},
		{U13_19, U13_20},
		{U13_21, U13_22},
		{U13_23, U13_24},
		{U13_25, U13_26},
		{U13_27, U13_28},
		{U13_29, U13_30},
		{U13_31, U13_32},
		{U13_33, U13_34},
		{U13_35, U13_36},
		{U13_37, U13_38},
		{U13_39, U13_40},
	})

	headers.Register("U14", [][]pins.Pin{
		{U14_1, U14_2},
		{U14_3, U14_4},
		{U14_5, U14_6},
		{U14_7, U14_8},
		{U14_9, U14_10},
		{U14_11, U14_12},
		{U14_13, U14_14},
		{U14_15, U14_16},
		{U14_17, U14_18},
		{U14_19, U14_20},
		{U14_21, U14_22},
		{U14_23, U14_24},
		{U14_25, U14_26},
		{U14_27, U14_28},
		{U14_29, U14_30},
		{U14_31, U14_32},
		{U14_33, U14_34},
		{U14_35, U14_36},
		{U14_37, U14_38},
		{U14_39, U14_40},
	})

	return true, nil
}

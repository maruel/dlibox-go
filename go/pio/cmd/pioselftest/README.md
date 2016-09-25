# pioselftest

Verifies that the library physically work. It requires the user to connect two
GPIO pins together and provide their pin number at the command line.

Sample output running on a Raspberry Pi:

```
$ pioselftest 12 6
Using drivers:
  - bcm283x
  - rpi
  - sysfs-gpio
  - sysfs-spi
  - sysfs-i2c
Using pins and their current state:
- GPIO12: In/High
- GPIO6: In/High

Testing GPIO6 -> GPIO12
  Testing base functionality
    GPIO12.In(Float)
    GPIO6.Out(Low)
    -> GPIO12: In/Low
    -> GPIO6: Out/Low
    GPIO6.Out(High)
    -> GPIO12: In/High
    -> GPIO6: Out/High
  Testing edges
    GPIO12.Edges()
    GPIO6.Out(Low)
    Low <- GPIO12
    GPIO6.Out(High)
    High <- GPIO12
    GPIO6.Out(Low)
    Low <- GPIO12
    GPIO12.DisableEdges()
  Testing pull resistor
    GPIO6.In(Down)
    -> GPIO12: In/Low
    -> GPIO6: In/Low
    GPIO6.In(Up)
    -> GPIO12: In/High
    -> GPIO6: In/High
Testing GPIO12 -> GPIO6
  Testing base functionality
    GPIO6.In(Float)
    GPIO12.Out(Low)
    -> GPIO6: In/Low
    -> GPIO12: Out/Low
    GPIO12.Out(High)
    -> GPIO6: In/High
    -> GPIO12: Out/High
  Testing edges
    GPIO6.Edges()
    GPIO12.Out(Low)
    Low <- GPIO6
    GPIO12.Out(High)
    High <- GPIO6
    GPIO12.Out(Low)
    Low <- GPIO6
    GPIO6.DisableEdges()
  Testing pull resistor
    GPIO12.In(Down)
    -> GPIO6: In/Low
    -> GPIO12: In/Low
    GPIO12.In(Up)
    -> GPIO6: In/High
    -> GPIO12: In/High
```

Sample output running on CHIP, note that pins 34 and 35 are the only combination that
supports all the functions being tested!

```
$ sudo ./pioselftest 35 34
Using drivers:
  - allwinner
  - chip
  - sysfs-gpio
  - sysfs-spi
  - sysfs-i2c
Using pins and their current state:
- PB3(35): In/High/Float
- PB2(34): In/High/Float

Testing PB2(34) -> PB3(35)
  Testing base functionality
    PB3(35).In(Float)
    PB2(34).Out(Low)
    -> PB3(35): In/Low/Float
    -> PB2(34): Out/Low
    PB2(34).Out(High)
    -> PB3(35): In/High/Float
    -> PB2(34): Out/High
  Testing edges
    PB3(35).Edges()
    PB2(34).Out(Low)
    Low <- PB3(35)
    PB2(34).Out(High)
    High <- PB3(35)
    PB2(34).Out(Low)
    Low <- PB3(35)
    PB3(35).DisableEdges()
  Testing pull resistor
    PB2(34).In(Down)
    -> PB3(35): EINT17
    -> PB2(34): In/Low/Down
    PB2(34).In(Up)
    -> PB3(35): EINT17
    -> PB2(34): In/High/Up
Testing PB3(35) -> PB2(34)
  Testing base functionality
    PB2(34).In(Float)
    PB3(35).Out(Low)
    -> PB2(34): In/Low/Float
    -> PB3(35): Out/Low
    PB3(35).Out(High)
    -> PB2(34): In/High/Float
    -> PB3(35): Out/High
  Testing edges
    PB2(34).Edges()
    PB3(35).Out(Low)
    Low <- PB2(34)
    PB3(35).Out(High)
    High <- PB2(34)
    PB3(35).Out(Low)
    Low <- PB2(34)
    PB2(34).DisableEdges()
  Testing pull resistor
    PB3(35).In(Down)
    -> PB2(34): EINT16
    -> PB3(35): In/Low/Down
    PB3(35).In(Up)
    -> PB2(34): EINT16
    -> PB3(35): In/High/Up
```

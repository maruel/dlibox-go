# pin-perf

Benchmarks sysfs vs direct pin access for gpio

Sample output on CHIP:
```
> sudo ./pin-perf 34 35
Using drivers:
  - allwinner
  - chip
  - sysfs-gpio
  - sysfs-spi
  - sysfs-i2c
Using pins and their current state:
- PB2(34)/GPIO34: Out/Low/Out/Low
- PB3(35)/GPIO35: In/Low/Float/In/Low

Testing gpio output speed
  1.2us per Out()
Testing sysfs output speed
  14.4us per Out()
Testing gpio input speed
  0.2us per In()
Testing sysfs input speed
  9.3us per In()
Testing gpio out+in speed
  1.1us per In()
Testing sysfs out+in speed
  26.7us per In()
```

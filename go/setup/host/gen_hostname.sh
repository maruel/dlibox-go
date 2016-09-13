#!/bin/bash
# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Calculate new hostname.
#
# Since the hostname is based on the serial number of the CPU with leading zeros
# trimmed off, it is a constant yet unique value.

set -eu

SERIAL="$(cat /proc/cpuinfo | grep Serial | cut -d ':' -f 2 | sed 's/^ 0\+//')"
if [ "$SERIAL" == "" ]; then
  SERIAL="$(hostnamectl status | grep 'Machine ID' | cut -d ':' -f 2 | cut -c 2- | cut -c -8)"
fi
echo "dlibox-$SERIAL"
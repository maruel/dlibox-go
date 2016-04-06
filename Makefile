# Copyright 2016 Marc-Antoine Ruel. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

# Set this variable to your host to enable "make push".
remote_host = raspberrypi2

# Use a trick to preinstall all imported packages. 'go build' doesn't permit
# installing packages, only 'go install' or 'go test -i' can do. But 'go
# install' would install an ARM binary, which is not what we want.
#
# Luckily, 'go test -i' is super fast on second execution.
dotstar: *.go cmd/dotstar/*.go
	GOOS=linux GOARCH=arm go test -i ./cmd/dotstar
	GOOS=linux GOARCH=arm go build ./cmd/dotstar

# When an executable is running, it must be scp'ed aside then moved over.
# dotstar will exit safely when it detects its binary changed so simple script
# like this works:
#
#     #!/bin/sh
#     set -e
#     while true; do
#       ./dotstar -port 8010
#     done
push: dotstar
	scp -q dotstar $(remote_host):dotstar2
	ssh $(remote_host) "mv dotstar2 dotstar"


# Runs it locally as a fake display with the web server running on port 8010.
run: *.go cmd/dotstar/*.go
	go install ./cmd/dotstar
	dotstar -fake -n 80 -port 8010


# Defaults to cross building to ARM.
all: dotstar

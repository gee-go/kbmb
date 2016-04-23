#!/bin/bash   
docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app iron/go:dev sh -c ls
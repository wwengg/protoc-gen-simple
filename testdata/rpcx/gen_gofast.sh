#!/bin/sh

protoc -I. -I${GOPATH}/src \
  --gogofast_out=. --gogofast_opt=paths=source_relative \
  --simple_out=. --simple_opt=paths=source_relative helloworld.proto

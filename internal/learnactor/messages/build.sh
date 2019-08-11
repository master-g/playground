#!/usr/bin/env bash

protoc -I=. -I="${GOPATH}/src" --gogoslick_out=plugins=grpc:. msg.proto

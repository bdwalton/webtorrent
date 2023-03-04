#!/bin/bash

protoc --experimental_allow_proto3_optional --go_out="${GOPATH}/src" models/*proto

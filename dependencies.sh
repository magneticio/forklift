#!/usr/bin/env bash

export GOPRIVATE=github.com/magneticio
go get -v ./...
set -x
go get github.com/stretchr/testify

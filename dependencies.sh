#!/usr/bin/env bash

go get -v ./...
set -x
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go get github.com/hashicorp/vault
go get github.com/hashicorp/vault/helper/strutil

#!/usr/bin/env bash
export GOPRIVATE=github.com/magneticio

if [ "$1" != "indocker" ]; then
  go test ./...
else
  GOIMAGE="dockercore/golang-cross:1.12.3"
  docker run --rm -it -v $(pwd):/src -w /src $GOIMAGE sh -c '
  for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
      echo "Building $GOOS-$GOARCH"
      export GOOS=$GOOS
      export GOARCH=$GOARCH
      if [ "$GOOS" = "windows" ]; then
        go get -u github.com/spf13/cobra
      fi
      go test ./...
    done
  done
  '
fi

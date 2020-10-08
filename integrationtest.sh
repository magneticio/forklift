#!/usr/bin/env bash
export GOPRIVATE=github.com/magneticio

if [ "$1" == "local" ]; then
  docker-compose up -d
  sleep 10
  VAULT_CLIENT_TOKEN=$(curl localhost:8201/client-token)
  curl --fail -v -H "X-Vault-Token:${VAULT_CLIENT_TOKEN}" localhost:8200/v1/sys/mounts
  curl --fail -v -X POST -H "X-Vault-Token:${VAULT_CLIENT_TOKEN}" -d '{"type": "kv"}' localhost:8200/v1/sys/mounts/secret
  go test -v -tags=integration -run=Integration ./...
  docker-compose down
elif [ "$1" == "circleci" ]; then
  sleep 10
  go test -v -tags=integration -run=Integration ./...
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
      go test -v -tags=integration -run=Integration ./...
    done
  done
  '
fi

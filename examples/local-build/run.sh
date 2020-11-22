#!/usr/bin/env bash

set -eo pipefail

if [[ $1 == "--clean" ]] ; then
  rm -rf ./libjq* ./example
  exit 0
fi

echo "=============================================="
echo "  Download and unpack pre-built libjq..."
echo "=============================================="
wget https://github.com/flant/libjq-go/releases/download/jq-b6be13d5-0/libjq-glibc-amd64.tgz
tar zxf libjq-glibc-amd64.tgz

echo "=============================================="
echo "  Build example.go"
echo "=============================================="

CGO_ENABLED=1 \
CGO_CFLAGS="-I$(pwd)/libjq/include" \
CGO_LDFLAGS="-L$(pwd)/libjq/lib" \
go build example.go

echo "=============================================="
echo "  Run example"
echo "=============================================="

./example

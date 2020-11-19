#!/usr/bin/env bash

set -eo pipefail

if [[ $1 == "--clean" ]] ; then
  docker rmi libjq-go-static:latest libjq-go-static:alpine libjq-go-static:debian
  exit 0
fi

echo "=============================================="
echo "  Build static binary"
echo "=============================================="
docker build . -t libjq-go-static:latest -f Dockerfile-static-binary

echo "=============================================="
echo "  Build alpine:3.7 final image"
echo "=============================================="
docker build . -t libjq-go-static:alpine -f Dockerfile-alpine3.7

echo "=============================================="
echo "  Build debian:jessie final image"
echo "=============================================="
docker build . -t libjq-go-static:debian -f Dockerfile-jessie

echo "=============================================="
echo "  Run example in scratch image"
echo "=============================================="
docker run --rm -ti libjq-go-static:latest

echo "=============================================="
echo "  Run example in alpine:3.7 image"
echo "=============================================="
docker run --rm -ti libjq-go-static:alpine

echo "=============================================="
echo "  Run example in debian:jessie image"
echo "=============================================="
docker run --rm -ti libjq-go-static:debian
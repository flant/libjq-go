#!/usr/bin/env bash

set -eo pipefail

if [[ $1 == "--clean" ]] ; then
  docker rmi libjq-go-simple:alpine libjq-go-simple:buster
  exit 0
fi

echo "=============================================="
echo "  Build debian:buster based image"
echo "=============================================="
docker build . -t libjq-go-simple:buster -f Dockerfile-buster
echo "=============================================="
echo "  Run example in debian:buster based image"
echo "=============================================="
docker run --rm -ti libjq-go-simple:buster -f Dockerfile-buster

echo "=============================================="
echo "  Build alpine:3.12 based image"
echo "=============================================="
docker build . -t libjq-go-simple:alpine -f Dockerfile-alpine
echo "=============================================="
echo "  Run example in alpine:3.12 based image"
echo "=============================================="
docker run --rm -ti libjq-go-simple:alpine -f Dockerfile-alpine

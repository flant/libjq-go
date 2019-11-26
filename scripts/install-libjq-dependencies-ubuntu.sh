#!/bin/sh

apt-get update

apt-get install -y \
    build-essential \
    autoconf \
    automake \
    libtool \
    git \
    bison \
    flex \
    wget

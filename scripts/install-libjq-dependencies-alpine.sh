#!/bin/sh

apk update

apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        git \
        autoconf \
        automake

apk add --virtual jq-deps \
        bison \
        flex \
        libtool

apk add bash



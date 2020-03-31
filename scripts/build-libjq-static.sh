#!/bin/sh

root=$1
out=$2

mkdir -p ${out}/libjq

echo "Build jq library into ${out}/libjq"

cd ${root}/modules/jq

git submodules update --init

autoreconf -fi
./configure CFLAGS=-fPIC --disable-maintainer-mode \
            --enable-all-static \
            --disable-shared \
            --disable-docs \
            --disable-valgrind \
            --with-oniguruma=builtin \
            --prefix=${out}/libjq
make
make install-libLTLIBRARIES install-includeHEADERS

echo Copy libonig

cp modules/oniguruma/src/.libs/libonig.a ${out}/libjq/lib
cp modules/oniguruma/src/.libs/libonig.la ${out}/libjq/lib
cp modules/oniguruma/src/.libs/libonig.lai ${out}/libjq/lib

echo "Use these flags with go build:"
echo "CGO_CFLAGS=-I${out}/libjq/include"
echo "CGO_LDFLAGS=-L${out}/libjq/lib"

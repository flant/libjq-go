#!/bin/sh

root=$1
out=$2

mkdir -p ${out}

echo "Build jq library into '${out}'"

cd ${root}/modules/jq

git submodule update --init

autoreconf -fi
./configure CFLAGS=-fPIC --disable-maintainer-mode \
            --enable-all-static \
            --disable-shared \
            --disable-docs \
            --disable-valgrind \
            --with-oniguruma=builtin \
            --prefix=${out}

make
make install-libLTLIBRARIES install-includeHEADERS

echo Copy libonig

cp modules/oniguruma/src/.libs/libonig.a ${out}/lib
cp modules/oniguruma/src/.libs/libonig.la ${out}/lib
cp modules/oniguruma/src/.libs/libonig.lai ${out}/lib

echo "Use these flags with go build:"
echo "CGO_CFLAGS=-I${out}/include"
echo "CGO_LDFLAGS=-L${out}/lib"

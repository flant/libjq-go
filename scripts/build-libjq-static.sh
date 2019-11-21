#!/bin/sh

root=$1
out=$2

mkdir -p ${out}/build/jq ${out}/build/onig

echo "Build oniguruma library into ${out}/build/oniguruma"

cd ${root}/modules/oniguruma

autoreconf -fi
./configure CFLAGS=-fPIC --disable-shared --prefix ${out}/build/onig
make
make install


echo "Build jq library into ${out}/build/jq"

cd ${root}/modules/jq

autoreconf -fi
./configure CFLAGS=-fPIC --disable-maintainer-mode \
            --enable-all-static \
            --disable-shared \
            --disable-docs \
            --disable-valgrind \
            --with-oniguruma=${out}/build/onig \
            --prefix=${out}/build/jq
make
make install-libLTLIBRARIES install-includeHEADERS

echo "Use these flags with go build:"
echo "CGO_CFLAGS=-I${out}/build/jq/include"
echo "CGO_LDFLAGS=\"-L${out}/build/onig/lib -L${out}/build/jq/lib\""

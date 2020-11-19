#!/usr/bin/env bash

# $1 is a jq repository path
# $2 is a output directory path

function usage() {
  cat <<EOF
Usage: $0 jq-path out-path

Example:
  $0 /app/jq /app/out
EOF
}

JQ_PATH=$1
if [[ $JQ_PATH == "" ]] ; then
  usage
  exit 1
fi

OUT=$2
if [[ $OUT == "" ]] ; then
  usage
  exit 1
fi

shaCmd="sha256sum"
if [[ "$OS" == "darwin" || "$OSTYPE" == "darwin"* ]] ; then
    shaCmd="shasum -a 256"
fi

cd $JQ_PATH

echo Build from commit $(git describe --tags)

autoreconf -fi

./configure CFLAGS=-fPIC --disable-maintainer-mode \
    --enable-all-static \
    --disable-shared \
    --disable-docs \
    --disable-tls \
    --disable-valgrind \
    --with-oniguruma=builtin --prefix=$OUT

# build jq and libjq
make -j2

# copy libjq.a and headers
make install-libLTLIBRARIES install-includeHEADERS
# copy libonig.a and headers
cp modules/oniguruma/src/.libs/libonig.* $OUT/lib
# copy bin/jq and docs
make install

strip $OUT/bin/jq

# Report jq version and a test.
echo "===================================="
echo "   JQ"
echo "------------------------------------"
echo -n "  version: "
$OUT/bin/jq --version
echo "------------------------------------"
echo "   Quick jq binary test:"
echo "------------------------------------"
echo '{"Key0":"asd", "NUMS":[1, 2, 3, 4, 5.123123123e+12]}' | $OUT/bin/jq '.,"Test OK"'
echo "===================================="
echo


# Generate checksum files.

mkdir $OUT/libjq
mv $OUT/lib $OUT/include $OUT/libjq
cd $OUT
find . -type f -exec sha256sum {} \; > $OUT/all.sha
cd $OUT/libjq
find . -type f -exec sha256sum {} \; > $OUT/libjq.sha

echo "===================================="
echo "   all.sha content:"
echo "===================================="
cat $OUT/all.sha
echo "===================================="
echo "   libjq.sha content:"
echo "===================================="
cat $OUT/libjq.sha
echo "===================================="
echo "   files in $OUT:"
echo "===================================="
find $OUT -exec ls -lad {} \;
echo "===================================="

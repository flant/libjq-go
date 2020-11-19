#!/usr/bin/env bash

function usage() {
  cat <<EOF
Extracts libjq files from /libjq directory in docker image.

Usage: $0 image:tag /directory/on/host

Example:
  $0 jq:b6be13d5 /app/libjq
EOF
}

IMAGE=$1
if [[ $IMAGE == "" ]] ; then
  usage
  exit 1
fi

OUT=$2
if [[ $OUT == "" ]] ; then
  usage
  exit 1
fi

mkdir -p $OUT/lib/pkgconfig $OUT/include
container_id=$(docker create $IMAGE)  # returns container ID
docker cp $container_id:/libjq/lib/libjq.la $OUT/lib
docker cp $container_id:/libjq/lib/libjq.a $OUT/lib
docker cp $container_id:/libjq/lib/libonig.la $OUT/lib
docker cp $container_id:/libjq/lib/libonig.a $OUT/lib
docker cp $container_id:/libjq/lib/libonig.lai $OUT/lib
docker cp $container_id:/libjq/lib/pkgconfig/libjq.pc $OUT/lib/pkgconfig
docker cp $container_id:/libjq/lib/pkgconfig/oniguruma.pc $OUT/lib/pkgconfig
docker cp $container_id:/libjq/include/oniguruma.h $OUT/include
docker cp $container_id:/libjq/include/onigposix.h $OUT/include
docker cp $container_id:/libjq/include/jv.h $OUT/include
docker cp $container_id:/libjq/include/oniggnu.h $OUT/include
docker cp $container_id:/libjq/include/jq.h $OUT/include
docker rm $container_id

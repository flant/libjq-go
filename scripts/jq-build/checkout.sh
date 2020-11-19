#!/usr/bin/env bash

# $1 is a commit  (b6be13d5)
# $2 is a target path (/jq)

JQ_GIT_SHA=$1
if [[ $JQ_GIT_SHA == "" ]] ; then
  usage
  exit 1
fi

JQ_PATH=$2
if [[ $JQ_PATH == "" ]] ; then
  usage
  exit 1
fi

git clone https://github.com/stedolan/jq.git $JQ_PATH
cd $JQ_PATH
git reset --hard $JQ_GIT_SHA
git submodule update --init

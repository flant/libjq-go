Scripts for .github/workflows/jq-build.yaml

# jq image release steps

1. Choose a commit in stedolan/jq repo.
2. Create tag in a form of  `jq[-<build_id>]-<commit>[-<serial>]`
3. Wait...

# Compatibility

jq static library built on buster (debian:10) is compatible with stretch (debian:9), ubuntu-18.04 and ubuntu 20.04.

alpine build is compatible with alpine 3.9+ — it was successfully tested in golang:1.11-alpine3.9 image.
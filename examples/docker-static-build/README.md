# Static build

You can build static binaries with libjq-go. For example, you can use image `golang:1.15-alpine` to build static binary and then COPY this binary to alpine or ubuntu or debian images. It will works even in older versions (jessie or alpine-3.7) and in a scratch image!

This example is an illustration of this setup.

## run it!

You need docker for this example.

Just execute `./run.sh` and it will build static binary in scratch image, create alpine and debian images and tests them all.

To clean up execute `./run.sh --clean`

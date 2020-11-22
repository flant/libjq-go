<p align="center">
<a href="https://hub.docker.com/r/flant/jq"><img src="https://img.shields.io/docker/pulls/flant/jq.svg?logo=docker" alt="docker pull flant/jq"/></a>
<a href="https://community.flant.com/c/shell-operator/7"><img src="https://img.shields.io/discourse/status?server=https%3A%2F%2Fcommunity.flant.com" alt="Discourse forum for discussions"/></a>
<a href="https://t.me/kubeoperator"><img src="https://img.shields.io/badge/telegram-RU%20chat-179cde.svg?logo=telegram" alt="Telegram chat RU"/></a>
</p>

# libjq-go

CGO bindings for jq with cache for compiled programs and ready-to-use static builds of libjq.

## Usage

```
import (
  "fmt"
  . "github.com/flant/libjq-go" // import Jq() shortcut
)

func main() {
	// 1. Run one jq program with one input.
	res, err = Jq().Program(".foo").Run(`{"foo":"bar"}`)

	// 2. Use directory with jq modules.
	res, err = Jq().WithLibPath("./jq_lib").
		Program(`....`).
		Run(`...`)
	
	// 3. Use program text as a key for a cache.
	for _, data := range inputJsons {
		res, err = Jq().Program(".foo").Cached().Run(data)
		// Do something with result ...
	}
	
	// 4. Explicitly precompile jq expression to speed up processing of multiple inputs.
	prg, err := Jq().Program(".foo").Precompile()
	for _, data := range inputJsons {
		res, err = prg.Run(data)
		// Do something with result ...
	}
	
	// 5. It is safe to use Jq() from multiple go-routines.
	//    Note however that programs are executed synchronously.
	go func() {
		res, err = Jq().Program(".foo").Run(`{"foo":"bar"}`)
	}()
	go func() {
		res, err = Jq().Program(".foo").Cached().Run(`{"foo":"bar"}`)
	}()
}
```

The full code is available in [example.go](./examples/simple/example.go).


# Build

## 1. Local build

The recommended and fastest way to build your program with libjq-go is to use pre-built static libjq, available in [Releases](https://github.com/flant/libjq-go/releases). See [local-build](./examples/local-build) example for inspiration.

```
wget https://github.com/flant/libjq-go/releases/download/jq-b6be13d5-0/libjq-glibc-amd64.tgz
tar zxf libjq-glibc-amd64.tgz
CGO_ENABLED=1 \
CGO_CFLAGS="-I./libjq/include" \
CGO_LDFLAGS="-L./libjq/lib" \
go build example.go
```

Also, you can use libjq in a form of dynamic library available in your OS or build static library from jq sources. Either way, read below about jq sources and performance issues.

### MacOS with brewed jq

```
brew install jq
CGO_ENABLED=1 go build example.go
```

## 2. Docker build

The recommended way for docker build is to use static libjq from [flant/jq](https://hub.docker.com/repository/docker/flant/jq) image available on hub.docker.com.

```
FROM flant/jq:b6be13d5-musl as libjq

FROM golang:1.15-alpine as builder
...
COPY --from=libjq /libjq /app/libjq/
...
RUN ... go build example.go

# Final image
FROM alpine:3.12
COPY --from=builder ...
```

Full source is available in [simple](./examples/simple) example.

If pre-built libjq is not an option, you can build static libjq in a separate image and then copy libjq to 'go builder' image. See this approach in a [Dockerfile](https://github.com/flant/shell-operator/blob/v1.0.0-beta.13/Dockerfile) of `flant/shell-operator` project.

## 3. Static build

Go can produce static binaries with CGO enabled. You should use static build of libjq and add -ldflags to `go build` command.

```
FROM flant/jq:b6be13d5-musl as libjq

FROM golang:1.15-alpine as builder
...
COPY --from=libjq /libjq /app/libjq/
...
RUN CGO_ENABLED=1 ... \
    go build \
    -ldflags="-linkmode external -extldflags '-static' -s -w" \
    example.go

# Final image
FROM alpine:3.12
COPY --from=builder ...
```

See [docker-static-build](./examples/docker-static-build) example.

# Notes

## jq source compatibility and jq 1.6 performance

TL;DR

- If your program works as a cli filter for one jq expression (like `jq` command itself) you can use any commit from `stedolan/jq`.
- If your program works as a server and process many jq expressions, consider use b6be13d5 commit and pre-built libjq assets.

Long story:

This library was tested with jq-1.5, jq-1.6 and with some commits from master branch. The submodule `jq` in this repository points to unreleased commit [stedolan/jq@b6be13d5](https://github.com/stedolan/jq/commit/b6be13d5de6dd7d8aad5fd871eb6b0b30fc7d7f6).

Which commit should you choose? Take these considerations into account:

- jq-1.5 works good, but it lucks new features.
- jq-1.6 turns out to be slow, see: [stedolan/jq#2069](https://github.com/stedolan/jq/issues/2069) and [flant/libjq-go#10](https://github.com/flant/libjq-go/issues/10).
- latest master have problem with `fromjson` and `tonumber` [stedolan/jq#2091](https://github.com/stedolan/jq/issues/2091).
- [stedolan/jq@b6be13d5](https://github.com/stedolan/jq/commit/b6be13d5de6dd7d8aad5fd871eb6b0b30fc7d7f6) is a commit that has features of v1.6 and a good performance and correctly handles errors in `fromjson`.


## Go compatibility

libjq-go is known to work with Go 1.11 and later versions.

## pre-built jq

For faster builds we publish pre-built libjq in [flant/jq](https://hub.docker.com/repository/docker/flant/jq) repository on hub.docker.com. Also, there are assets in jq-* [releases](https://github.com/flant/libjq-go/releases) in this repo on github.

`glibc` build is known to work in debian:stretch, debian:buster, ubuntu:18.04, ubuntu:20.04 (and seems to work in alpine).

`musl` build is known to work in alpine:3.7+ (You can even compile your program using `golang:1.11-alpine3.7`!)

> Note: `flant/jq` image contains /bin/jq binary.

## Inspired projects

There are other `jq` bindings in Go:

- https://github.com/aki017/gq
- https://github.com/bongole/go-jq
- https://github.com/mgood/go-jq
- https://github.com/threatgrid/jq-go
- https://github.com/mattatcha/jq
- https://github.com/jzelinskie/faq

Also these projects was very helpful in understanding jq sources:

- https://github.com/robertaboukhalil/jqkungfu
- https://github.com/doloopwhile/pyjq


## License

Apache License 2.0, see [LICENSE](LICENSE).

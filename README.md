# libjq-go

CGO bindings for jq with cache for compiled programs.

## Usage

```
import (
  "fmt"
  . "github.com/flant/libjq-go" // import Jq() shortcut
)

func main() {
	// 1. Run one program with one input.
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
	
	// 4. Explicitly precompile jq expression.
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

This code is available in [example.go](example/example.go) as a working example.


## Build

1. Local build

To build your program with this library, you should install some build dependencies and statically compile oniguruma and jq libraries:

```
apt-get install build-essential autoconf automake libtool flex bison wget
git clone https://github.com/flant/libjq-go
cd libjq-go
git submodule update --init
./scripts/build-libjq-static.sh $(pwd) $(pwd)/out
export LIBJQ_CFLAGS="-I$(pwd)/out/jq/include"
export LIBJQ_LDFLAGS="-L$(pwd)/out/oniguruma/lib -L$(pwd)/out/jq/lib"
```

Now you can build your application:

```
CGO_ENABLED=1 CGO_CFLAGS="${LIBJQ_CFLAGS}" CGO_LDFLAGS="${LIBJQ_LDFLAGS}" go build <your arguments>
```

2. Docker build

If you want to build your program with docker, you can build oniguruma and jq in artifact image and then copy them to go builder image. See example of this approach in [Dockerfile](https://github.com/flant/shell-operator/blob/master/Dockerfile) of a shell-operator — the real project that use this library.


## jq source compatibility and jq 1.6 performance

This library was tested with jq-1.5, jq-1.6 and with some commits from master branch. The submodule `jq` in this repository points to unreleased commit [stedolan/jq@b6be13d](https://github.com/stedolan/jq/commit/b6be13d5de6dd7d8aad5fd871eb6b0b30fc7d7f6).

Which commit should you choose? Take these considerations into account:

- jq-1.5 works good, but it lucks new features.
- jq-1.6 turns out to be slow, see: [stedolan/jq#2069](https://github.com/stedolan/jq/issues/2069) and #10.
- latest master have problem with `fromjson` and `tonumber` [stedolan/jq#2091](https://github.com/stedolan/jq/issues/2091).
- [stedolan/jq@b6be13d](https://github.com/stedolan/jq/commit/b6be13d5de6dd7d8aad5fd871eb6b0b30fc7d7f6) is a commit that is later than jq-1.6, works fast and correctly handles errors in `fromjson`.


## Inspired projects

There are other `jq` bindings in Go:

- https://github.com/aki017/gq
- https://github.com/bongole/go-jq
- https://github.com/mgood/go-jq
- https://github.com/threatgrid/jq-go
- https://github.com/mattatcha/jq
- https://github.com/jzelinskie/faq

Also these projects was very helpful: 

- https://github.com/robertaboukhalil/jqkungfu
- https://github.com/doloopwhile/pyjq


## License

Apache License 2.0, see [LICENSE](LICENSE).

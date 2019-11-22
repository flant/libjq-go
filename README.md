# libjq-go

CGO bindings for jq with cache for compiled programs

## Usage

```
import (
  "fmt"
  . "github.com/flant/libjq-go" // import Jq, JqMainThread and JqCallLoop
)

func main() {
  // Jq instance with direct calls of libjq methods. Note that it cannot be used in go routines.
  var jq = JqMainThread

  // Run one program with one input.
  res, err := jq().Program(".foo").Run(`{"foo":"bar"}`)

  // Use directory with jq modules.
  res, err := jq().WithLibPath("./jq_lib").
    Program(...).
    Run(...)

  // Use jq state cache to speedup handling of multiple inputs.
  prg, err := jq().Program(...).Precompile()
  for _, data := range InputJsons {
    res, err = prg.Run(data)
    // do something with filter result ...
  }

  // Use jq from go-routines.
  // Jq() helper returns instance that use LockOsThread trick to run libjq methods in main thread.
  done := make(chan struct{})

  go func() {
    res, err := Jq().Program(".foo").Run(`{"foo":"bar"}`)
    done <- struct{}{}
  }()

  // main is locked here.
  JqCallLoop(done)
}
```


## Build

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

# libjq-go
CGO bindings for jq with cache for compiled programs

## Example

```
import (
  "fmt"
  . "github.com/flant/libjq-go"
)

func main() {
  // Run one program with one input.
  res, err := Jq().Program(".foo").Run(`{"foo":"bar"}`)
  if err != nil {
    panic(err)
  }
  fmt.Printf("filter result: %s\n", res)
  
  // Use jq state cache to speedup handling of multiple inputs.
  jqp, err := Jq().Program(".[]|.bar").Compile()
  if err != nil {
    panic(err)
  }
  for _, data := range InputJson {
    res, err := jqp.Run(data)
    // do something with filter result ...
  }
  
  // Use library directory.
  JqWithLib("./jq_lib").Program(`include "libname"; .foo|libmethod`).Run(`{"foo":"json here"}`)
  
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

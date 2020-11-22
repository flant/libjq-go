package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	. "github.com/flant/libjq-go"
)

/*
Run it locally if oniguruma and jq are compiled:

CGO_ENABLED=1 \
CGO_CFLAGS="-Ipath-to-jq_lib/include" \
CGO_LDFLAGS="-Lpath-to-oniguruma_lib/lib -Lpath-to-jq_lib/lib" \
go run example.go

1. "bar"
2. kebab-string-here "kebabStringHere"
3. "bar-quux"
3. "baz-baz"
4. "Foo quux"
4. "Foo baz"
5. "bar"
5. "bar"
5. "bar"

*/

func main() {
	var res string
	var err error
	var inputJsons []string

	// 1. Run one program with one input.
	res, err = Jq().Program(".foo").Run(`{"foo":"bar"}`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("1. %s\n", res)
	// Should print
	// 1. "bar"

	// 2. Use directory with jq modules.
	prepareJqLib()
	res, err = Jq().WithLibPath("./jq_lib").
		Program(`include "mylibrary"; .foo|mymethod`).
		Run(`{"foo":"kebab-string-here"}`)
	fmt.Printf("2. %s %s\n", "kebab-string-here", res)
	removeJqLib()
	// Should print
	// 2. kebab-string-here "kebabStringHere"

	// 3. Use program text as a key for a cache.
	inputJsons = []string{
		`{ "foo":"bar-quux" }`,
		`{ "foo":"baz-baz" }`,
		// ...
	}
	for _, data := range inputJsons {
		res, err = Jq().Program(".foo").Cached().Run(data)
		if err != nil {
			panic(err)
		}
		// Now do something with filter result ...
		fmt.Printf("3. %s\n", res)
	}
	// Should print
	// 3. "bar-quux"
	// 3. "baz-baz"

	// 4. Explicitly precompile jq expression.
	inputJsons = []string{
		`{ "bar":"Foo quux" }`,
		`{ "bar":"Foo baz" }`,
		// ...
	}
	prg, err := Jq().Program(".bar").Precompile()
	if err != nil {
		panic(err)
	}
	for _, data := range inputJsons {
		res, err = prg.Run(data)
		if err != nil {
			panic(err)
		}
		// Now do something with filter result ...
		fmt.Printf("4. %s\n", res)
	}
	// Should print
	// 4. "Foo quux"
	// 4. "Foo baz"

	// 5. It is safe to use Jq() from multiple go-routines.
	//    Note however that programs are executed synchronously.
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		res, err = Jq().Program(".foo").Run(`{"foo":"bar"}`)
		if err != nil {
			panic(err)
		}
		fmt.Printf("5. %s\n", res)
		wg.Done()
	}()
	go func() {
		res, err = Jq().Program(".foo").Cached().Run(`{"foo":"bar"}`)
		if err != nil {
			panic(err)
		}
		fmt.Printf("5. %s\n", res)
		wg.Done()
	}()
	go func() {
		res, err = Jq().Program(".foo").Cached().Run(`{"foo":"bar"}`)
		if err != nil {
			panic(err)
		}
		fmt.Printf("5. %s\n", res)
		wg.Done()
	}()
	wg.Wait()
	// Should print
	// 5. "bar"
	// 5. "bar"
	// 5. "bar"
}

func prepareJqLib() {
	if err := os.MkdirAll("./jq_lib", 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("./jq_lib/mylibrary.jq", []byte(`def mymethod: gsub("-(?<a>[a-z])"; .a|ascii_upcase);`), 0644); err != nil {
		panic(err)
	}
}

func removeJqLib() {
	if err := os.RemoveAll("./jq_lib"); err != nil {
		panic(err)
	}
}

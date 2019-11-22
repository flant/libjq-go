package main

import (
	"fmt"

	. "github.com/flant/libjq-go"
)

func main() {
	var res string
	var err error

	// Jq instance with direct calls of libjq methods — cannot be used is go routines.
	var jq = JqMainThread

	// Run one program with one input.
	res, err = jq().Program(".foo").Run(`{"foo":"bar"}`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("filter result: %s\n", res)

	// Use jq state cache to speedup handling of multiple inputs.
	InputJson := []string{
		`{..fields..}`,
		`{..fields..}`,
	}
	jqp, err := jq().Program(".[]|.bar").Precompile()
	if err != nil {
		panic(err)
	}
	for _, data := range InputJson {
		res, err = jqp.Run(data)
		// do something with filter result ...
	}

	// Use directory with jq modules.
	res, err = jq().WithLibPath("./jq_lib").
		Program(`include "libname"; .foo|libmethod`).
		Run(`{"foo":"json here"}`)

	// Use jq from go-routines.
	// Jq() returns instance that use LockOsThread trick to run libjq methods in main thread.
	done := make(chan struct{})

	go func() {
		res, err = Jq().Program(".foo").Run(`{"foo":"bar"}`)
		done <- struct{}{}
	}()

	// main is locked here.
	JqCallLoop(done)
}

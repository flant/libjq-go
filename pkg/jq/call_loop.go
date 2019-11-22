package jq

import "runtime"

/*
libjq methods should run from main thread, so this trick with LockOsThread come up.
*/

func init() {
	runtime.LockOSThread()
}

// JqCallLoop handles external functions. This loop should be started from the main.main.
func JqCallLoop(done chan struct{}) {
	for {
		select {
		case f := <-jqcalls:
			f()
		case <-done:
			return
		}
	}
}

var jqcalls = make(chan func())

// JqCall is used to run jq related code in main thread.
func JqCall(f func()) {
	done := make(chan struct{}, 1)
	jqcalls <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

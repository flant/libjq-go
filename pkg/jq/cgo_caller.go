package jq

import (
	"runtime"
	"sync"
)

/*
libjq methods should run in one thread, so this trick with LockOsThread come up.
*/

var cgoCallsCh chan func()
var mu = sync.Mutex{}

// CgoCall is used to run C code of a jq in a dedicated go-routine locked to OS thread.
func CgoCall(f func()) {
	mu.Lock()
	if cgoCallsCh == nil {
		cgoCallsCh = make(chan func())
		go func() {
			runtime.LockOSThread()
			for {
				select {
				case f := <-cgoCallsCh:
					f()
				}
			}
		}()
	}
	mu.Unlock()

	done := make(chan struct{}, 1)
	cgoCallsCh <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

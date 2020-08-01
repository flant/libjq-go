package jq

import (
	"runtime"
	"sync"
)

/*
jq_state is thread safe, but not compatible with migration of go routines between thread.
That is why libjq methods should be called from the same thread where jq_state was created.
To achieve this, a dedicated go routine and a chan func() are used.
*/

type CgoCaller func(func())

// NewCgoCaller is a factory of CgoCallers. CgoCaller is a way to run C code of a jq in a dedicated go-routine locked to OS thread.
// CgoCaller on first invoke creates a channel and starts a go-routine locked to os thread. This go-routine receives tasks to run via a channel.

func NewCgoCaller() CgoCaller {
	var cgoCallTasksCh chan func()
	var initOnce sync.Once

	return func(f func()) {
		initOnce.Do(func() {
			cgoCallTasksCh = make(chan func())
			go func() {
				runtime.LockOSThread()
				for {
					select {
					case f := <-cgoCallTasksCh:
						f()
					}
				}
			}()
		})

		var wg sync.WaitGroup
		wg.Add(1)
		cgoCallTasksCh <- func() {
			f()
			wg.Done()
		}
		wg.Wait()
	}
}

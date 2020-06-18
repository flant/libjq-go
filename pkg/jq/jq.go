package jq

import (
	"github.com/flant/libjq-go/pkg/libjq"
)

/*
High level API for libjq.

Executor for jq programs has this options:
- cache
- library path
- cgo caller function to pin jq calls to an OS thread
*/
type Jq struct {
	Cache     *JqCache
	LibPath   string
	CgoCaller CgoCaller
}

func NewJq() *Jq {
	return &Jq{}
}

func (jq *Jq) WithCache(cache *JqCache) *Jq {
	jq.Cache = cache
	return jq
}

func (jq *Jq) WithCgoCaller(cgoCaller CgoCaller) *Jq {
	jq.CgoCaller = cgoCaller
	return jq
}

func (jq *Jq) WithLibPath(path string) *Jq {
	jq.LibPath = path
	return jq
}

func (jq *Jq) Program(program string) *JqProgram {
	return &JqProgram{
		Jq:      jq,
		Program: program,
	}
}

type JqProgram struct {
	Jq          *Jq
	Program     string
	CacheLookup bool
}

// Cached set cached flag so next call to Run will put compiled program to cache.
func (jqp *JqProgram) Cached() *JqProgram {
	if jqp.Jq.Cache != nil {
		jqp.CacheLookup = true
	}
	return jqp
}

// Precompile can be used to compile a program and store jq state in cache.
// Method returns error in case of syntax error.
func (jqp *JqProgram) Precompile() (p *JqProgram, err error) {
	if jqp.Jq.Cache == nil {
		return jqp, nil
	}

	jqp.CacheLookup = true

	_, err = jqp.compile()
	return jqp, err
}

// Compile compiles a program immediately, it returns error in case of syntax error.
func (jqp *JqProgram) Compile() (p *JqProgram, err error) {
	_, err = jqp.compile()
	return jqp, err
}

// Run actually runs a program over passed data. It compiles program
// if the program is not compiled yet.
// Returns an quoted string if filter result is a string.
func (jqp *JqProgram) Run(data string) (s string, e error) {
	return jqp.run(data, false)
}

// RunRaw actually runs a program over passed data. It compiles program
// if the program is not compiled yet.
// Returns an unquoted string if filter result is a string.
func (jqp *JqProgram) RunRaw(data string) (s string, e error) {
	return jqp.run(data, true)
}

// compile creates a new jq state with compiled program or just returns a cached one.
func (jqp *JqProgram) compile() (state *libjq.JqState, err error) {
	if jqp.CacheLookup {
		inCacheState := jqp.Jq.Cache.Get(jqp.Program)
		if inCacheState != nil {
			return inCacheState, nil
		}
	}

	jqp.Jq.CgoCaller(func() {
		state, err = libjq.NewJqState()
		if err != nil {
			return
		}
		state.SetLibraryPath(jqp.Jq.LibPath)
		err = state.Compile(jqp.Program)
	})
	if err != nil {
		return
	}

	if jqp.CacheLookup {
		jqp.Jq.Cache.Set(jqp.Program, state)
	}
	return state, nil
}

// run starts compiled program.
func (jqp *JqProgram) run(inJson string, rawMode bool) (res string, err error) {
	var state *libjq.JqState
	state, err = jqp.compile()
	if err != nil {
		return "", err
	}

	jqp.Jq.CgoCaller(func() {
		res, err = state.ProcessOneValue(inJson, rawMode)
		if !jqp.CacheLookup {
			state.Teardown()
		}
	})

	return
}

package jq

import (
	"github.com/flant/libjq-go/pkg/libjq"
)

/*
High level API for libjq.

Jq has this options:
- cache
- library path
*/
type Jq struct {
	Cache   *JqCache
	LibPath string
}

func NewJq() *Jq {
	return &Jq{}
}

func (jq *Jq) WithCache(cache *JqCache) *Jq {
	jq.Cache = cache
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

	CgoCall(func() {
		_, err = jqp.compile()
	})

	return jqp, err
}

// Run actually runs a program over passed data. It compiles program
// if the program is not compiled yet.
func (jqp *JqProgram) Run(data string) (s string, e error) {
	CgoCall(func() {
		s, e = jqp.run(data, false)
	})
	return
}

// RunRaw actually runs a program over passed data. It compiles program
// if the program is not compiled yet.
// Returns an unquoted string if filter result is a string.
func (jqp *JqProgram) RunRaw(data string) (s string, e error) {
	CgoCall(func() {
		s, e = jqp.run(data, true)
	})
	return
}

// compile create a jq state with compiled program and stores it in cache if needed.
func (jqp *JqProgram) compile() (state *libjq.JqState, err error) {
	if jqp.CacheLookup {
		inCacheState := jqp.Jq.Cache.Get(jqp.Program)
		if inCacheState != nil {
			return inCacheState, nil
		}
	}
	state, err = libjq.NewJqState()
	if err != nil {
		return nil, err
	}
	state.SetLibraryPath(jqp.Jq.LibPath)
	err = state.Compile(jqp.Program)
	if err != nil {
		return nil, err
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
	if !jqp.CacheLookup {
		defer state.Teardown()
	}

	return state.ProcessOneValue(inJson, rawMode)
}

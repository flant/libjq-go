package jq

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup JqCallLoop for jq testing
func TestMain(m *testing.M) {
	exitCode := 0
	done := make(chan struct{})

	go func() {
		exitCode = m.Run()
		done <- struct{}{}
	}()

	JqCallLoop(done)

	os.Exit(exitCode)
}

func Test_FieldAccess(t *testing.T) {
	res, err := NewJq().WithCallProxy(JqCall).
		Program(".foo").Run(`{"foo":"baz"}`)
	assert.NoError(t, err)
	assert.Equal(t, `"baz"`, res)
}

func Test_JsonOutput(t *testing.T) {
	in := `{"foo":"baz","bar":"quux"}`
	res, err := NewJq().WithCallProxy(JqCall).
		Program(".").Run(in)
	assert.NoError(t, err)
	assert.Equal(t, in, res)
}

func Test_LibPath_FilteredFieldAccess(t *testing.T) {
	prg := `include "camel"; .bar | camel`
	in := `{"foo":"baz","bar":"quux-mooz"}`
	out := `"quuxMooz"`

	res, err := NewJq().WithCallProxy(JqCall).
		WithLibPath("./testdata/jq_lib").
		Program(prg).Run(in)
	assert.NoError(t, err)
	assert.Equal(t, out, res)
}

func Test_CachedProgram_FieldAccess(t *testing.T) {
	p, err := NewJq().WithCallProxy(JqCall).
		WithCache(JqDefaultCache()).
		Program(".foo").Precompile()
	assert.NoError(t, err)

	for i := 0; i < 50; i++ {
		val := fmt.Sprintf(`"baz%d"`, i)
		in := fmt.Sprintf(`{"foo":%s}`, val)
		res, err := p.Run(in)
		assert.NoError(t, err)
		assert.Equal(t, val, res)
	}
}

func Test_Concurent_FieldAccess(t *testing.T) {
	job := func() {
		for i := 0; i < 50; i++ {
			prg := fmt.Sprintf(`include "camel"; .foo%d | camel`, i)
			val := fmt.Sprintf(`"quux-baz%d"`, i)
			out := fmt.Sprintf(`"quuxBaz%d"`, i)
			in := fmt.Sprintf(`{"foo%d":%s}`, i, val)

			res, err := NewJq().WithCallProxy(JqCall).
				WithCache(JqDefaultCache()).
				WithLibPath("./testdata/jq_lib").
				Program(prg).Cached().Run(in)
			assert.NoError(t, err)
			assert.Equal(t, out, res)
		}
	}

	parallelism := 16

	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			job()
			wg.Done()
		}()
	}
	wg.Wait()
}

// TODO catch errors: syntax, input, program run

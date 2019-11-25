package jq

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

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

	parallelism := 32

	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			if parallelism%2 == 0 {
				runtime.LockOSThread()
			}
			job()
			wg.Done()
		}()
	}
	wg.Wait()
}

// TODO tests to catch jq processing errors: syntax, input and program run

// run to catch memory leaks!
func Test_LongRunner_BigData(t *testing.T) {
	t.SkipNow()

	parallelism := 16

	// There are `parallelism` of different programs and fooXXX fields,
	// but extra field is always different.
	job := func(jobId int) {
		i := 100000
		for {
			prg := fmt.Sprintf(`include "camel"; .foo%d | camel`, i%parallelism)
			val := fmt.Sprintf(`"quux-baz%d"`, i%parallelism)
			out := fmt.Sprintf(`"quuxBaz%d"`, i%parallelism)
			in := fmt.Sprintf(`{"foo%d":%s, "extra":%s}`, i%parallelism, val, generateBigJsonObject(1024, i))

			res, err := NewJq().WithCallProxy(JqCall).
				WithCache(JqDefaultCache()).
				WithLibPath("./testdata/jq_lib").
				Program(prg).Cached().Run(in)
			if assert.NoError(t, err) {
				assert.Equal(t, out, res)
			}
			i--
			if i == 0 {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {

		go func(jobId int) {
			secs := time.Duration(2 * jobId)
			time.Sleep(secs * time.Second)
			fmt.Printf("Start %d\n", jobId)
			job(jobId)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func generateBigJsonObject(size int, id int) string {
	var b strings.Builder

	b.WriteString(`{"a":"`)

	bt := make([]byte, size)
	for i := 0; i < len(bt); i++ {
		bt[i] = ' '
	}
	// Put X somewher
	bt[id%(len(bt))] = 'X'

	b.Write(bt)
	b.WriteString(`"}`)
	return b.String()
}

func Test_BigObject(t *testing.T) {
	t.SkipNow()
	for i := 0; i < 100; i++ {
		fmt.Println(generateBigJsonObject(25, i))
	}
}

package jq

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func Test_FieldAccess(t *testing.T) {
	g := NewWithT(t)

	res, err := NewJq().Program(".foo").Run(`{"foo":"baz"}`)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(`"baz"`))
}

func Test_JsonOutput(t *testing.T) {
	g := NewWithT(t)
	in := `{"foo":"baz","bar":"quux"}`
	res, err := NewJq().Program(".").Run(in)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(in))
}

func Test_LibPath_FilteredFieldAccess(t *testing.T) {
	g := NewWithT(t)

	prg := `include "camel"; .bar | camel`
	in := `{"foo":"baz","bar":"quux-mooz"}`
	out := `"quuxMooz"`

	res, err := NewJq().WithLibPath("./testdata/jq_lib").
		Program(prg).Run(in)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(out))
}

func Test_CachedProgram_FieldAccess(t *testing.T) {
	g := NewWithT(t)

	p, err := NewJq().WithCache(JqDefaultCache()).
		Program(".foo").Precompile()
	g.Expect(err).ShouldNot(HaveOccurred())

	for i := 0; i < 50; i++ {
		val := fmt.Sprintf(`"baz%d"`, i)
		in := fmt.Sprintf(`{"foo":%s}`, val)
		res, err := p.Run(in)
		g.Expect(err).ShouldNot(HaveOccurred())
		g.Expect(res).To(Equal(val))
	}
}

func Test_Concurent_FieldAccess(t *testing.T) {
	g := NewWithT(t)

	job := func() {
		for i := 0; i < 50; i++ {
			prg := fmt.Sprintf(`include "camel"; .foo%d | camel`, i)
			val := fmt.Sprintf(`"quux-baz%d"`, i)
			out := fmt.Sprintf(`"quuxBaz%d"`, i)
			in := fmt.Sprintf(`{"foo%d":%s}`, i, val)

			res, err := NewJq().
				WithCache(JqDefaultCache()).
				WithLibPath("./testdata/jq_lib").
				Program(prg).Cached().Run(in)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(res).To(Equal(out))
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

// TODO add more tests to catch jq processing errors: syntax, input and program run

// Uncomment SkipNow to run and catch memory leaks!
// TODO add script to run test and watch for memory leaks
func Test_LongRunner_BigData(t *testing.T) {
	t.SkipNow()
	g := NewWithT(t)

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

			res, err := NewJq().
				WithCache(JqDefaultCache()).
				WithLibPath("./testdata/jq_lib").
				Program(prg).Cached().Run(in)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(res).To(Equal(out))

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
	g := NewWithT(t)

	g.Expect(generateBigJsonObject(25, 0)).To(Equal(`{"a":"X                        "}`))
	g.Expect(generateBigJsonObject(25, 9)).To(Equal(`{"a":"         X               "}`))
	g.Expect(generateBigJsonObject(25, 24)).To(Equal(`{"a":"                        X"}`))
	g.Expect(generateBigJsonObject(25, 25)).To(Equal(generateBigJsonObject(25, 0)))
	g.Expect(generateBigJsonObject(25, 49)).To(Equal(generateBigJsonObject(25, 24)))
}

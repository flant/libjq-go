package libjq_go

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/flant/libjq-go/pkg/jq"
)

func Test_OneProgram_OneInput(t *testing.T) {
	g := NewWithT(t)

	res, err := Jq().Program(".foo").Run(`{"foo":"bar"}`)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(`"bar"`))

	res, err = Jq().Program(".foo").RunRaw(`{"foo":"bar"}`)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(`bar`))
}

func Benchmark_HasKey(b *testing.B) {

	p, err := Jq().Program(`has("foo")`).Precompile()
	if err != nil {
		b.Fatalf("precompile program: %s", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := p.Run(`{"bar":"baz"}`)
		if err != nil {
			b.Fatalf("run %d: %s", i, err)
		}
	}
}

func Benchmark_PreCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prg := fmt.Sprintf(`has("foo%d")`, i)
		p, err := Jq().Program(prg).Precompile()
		if err != nil {
			b.Fatalf("precompile program %d: %s", i, err)
		}
		_, err = p.Run(`{"bar":"baz"}`)
		if err != nil {
			b.Fatalf("run %d: %s", i, err)
		}
	}
}

func Test_CompileError(t *testing.T) {
	g := NewWithT(t)

	_, err := Jq().Program(`{"message": .message"}`).Run(`{"message":"bar"}`)

	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("jq: error: syntax error"))
	g.Expect(err.Error()).To(ContainSubstring("compile error"))
	g.Expect(err.Error()).ToNot(ContainSubstring("0 0 0")) // {0 0 0 0 [0 0 0 0 0 0 0 0]}
}

func Test_RunError(t *testing.T) {
	g := NewWithT(t)
	_, err := Jq().Program(".foo[] | keys").Run(`{"foo":"bar"}`)

	g.Expect(err).Should(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("Cannot iterate over string"))
}

// PoC of multiple jq thread
func Test_two_cgo_callers(t *testing.T) {
	g := NewWithT(t)

	cgoCaller := jq.NewCgoCaller()
	invoker := jq.NewJq().
		WithCache(jq.NewJqCache()).
		WithCgoCaller(cgoCaller)
	// New cache is required for every new cgoCaller!
	// Cache saves jq_state from jq_compile in a pinned thread, so another cgoCaller
	// cannot access that jq_state in another thread.
	//WithCache(jq.JqDefaultCache()):
	// Assertion failed: (0 && "invalid instruction"), function jq_nextAssertion failed: (jv_is_valid(v, file src/execute.c, line 401.
	// al)), function stack_pop, file src/execute.c, line 177.
	// SIGABRT: abort

	var wg sync.WaitGroup

	wg.Add(2)

	// jq with default cgo caller
	go func() {
		defer wg.Done()

		p1, _ := Jq().Program(`.bb//"NO"`).Precompile()

		for i := 0; i < 500; i++ {
			res, err := Jq().Program(".foo").Run(`{"foo":"bar"}`)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(res).To(Equal(`"bar"`))

			res, err = p1.Run(`{"foo":"bar"}`)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(res).To(Equal(`"NO"`))
		}

	}()

	// jq with another, parallel cgo caller
	go func() {
		defer wg.Done()

		p1, _ := invoker.Program(`.bb//"NO"`).Precompile()

		for i := 0; i < 500; i++ {
			res, err := p1.Run(`{"foo":"bar"}`)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(res).To(Equal(`"NO"`))

			res, err = invoker.Program(".foo").Run(`{"foo":"bar"}`)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(res).To(Equal(`"bar"`))
		}

	}()

	wg.Wait()
}

// PoC of multiple jq thread
func Test_multiple_cgo_callers(t *testing.T) {
	g := NewWithT(t)

	jqPoolLen := 16

	// init invokers
	jqPool := []*jq.Jq{}
	for i := 0; i < jqPoolLen; i++ {
		invoker := jq.NewJq().
			WithCache(jq.NewJqCache()).
			WithCgoCaller(jq.NewCgoCaller())
		jqPool = append(jqPool, invoker)
	}

	consumersCount := jqPoolLen * 2 // twice as invokers

	// Start consumers
	var wg sync.WaitGroup
	wg.Add(consumersCount)

	for i := 0; i < consumersCount; i++ {
		invokerIndex := i % jqPoolLen
		go func(n int) {
			defer wg.Done()

			invoker := jqPool[n]

			p1, _ := invoker.Program(`.bb//"NO"`).Precompile()

			foo, _ := invoker.Program(".foo").Precompile()

			for i := 0; i < 100; i++ {
				res, err := p1.Run(`{"foo":"bar"}`)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(res).To(Equal(`"NO"`))

				res, err = foo.Run(`{"foo":"bar"}`)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(res).To(Equal(`"bar"`))
			}

			p2, _ := invoker.Program(`.bb//"YES"`).Precompile()

			p3, _ := invoker.Program(".foobar").Precompile()

			for i := 0; i < 100; i++ {
				res, err := p2.Run(`{"foo":"bar"}`)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(res).To(Equal(`"YES"`))

				res, err = p3.Run(`{"foobar":"bar"}`)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(res).To(Equal(`"bar"`))
			}

		}(invokerIndex)
	}

	wg.Wait()
}

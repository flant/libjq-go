package libjq_go

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
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

package libjq_go

import (
	"fmt"
	"strings"
	"testing"
)

func Test_OneProgram_OneInput(t *testing.T) {

	res, err := Jq().Program(".foo").Run(`{"foo":"bar"}`)
	if err != nil {
		t.Fatalf("expect program not fail: %s", err)
	}
	if res != `"bar"` {
		t.Fatalf("expect '\"bar\"', got '%s'", res)
	}

	res, err = Jq().Program(".foo").RunRaw(`{"foo":"bar"}`)
	if err != nil {
		t.Fatalf("expect program not fail: %s", err)
	}
	if res != `bar` {
		t.Fatalf("expect 'bar', got '%s'", res)
	}
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
	_, err := Jq().Program(`{"message": .message"}`).Run(`{"message":"bar"}`)

	if err == nil {
		t.Fatal("expect program should fail")
	}
	expect := "jq: error: syntax error"
	if !strings.Contains(err.Error(), expect) {
		t.Fatalf("expect '%s' in err: %s", expect, err)
	}
	expect = "compile error"
	if !strings.Contains(err.Error(), expect) {
		t.Fatalf("expect '%s' in err: %s", expect, err)
	}
	expect = "0 0 0" // {0 0 0 0 [0 0 0 0 0 0 0 0]} problem
	if strings.Contains(err.Error(), expect) {
		t.Fatalf("not expect '%s' in err: %s", expect, err)
	}
}

func Test_RunError(t *testing.T) {
	_, err := Jq().Program(".foo[] | keys").Run(`{"foo":"bar"}`)

	if err == nil {
		t.Fatal("expect program should fail")
	}
	expect := "Cannot iterate over string"
	if !strings.Contains(err.Error(), expect) {
		t.Fatalf("expect '%s' in err: %s", expect, err)
	}
}

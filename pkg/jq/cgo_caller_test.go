package jq

import (
	"github.com/flant/libjq-go/pkg/libjq"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_CgoCall(t *testing.T) {
	g := NewWithT(t)

	testProgram := `.foo`
	testData := `{"foo": "bar"}`
	testExpected := `"bar"`
	testResult := ""

	var jqState *libjq.JqState
	var err error

	caller := NewCgoCaller()

	// Create jq state in locked OS thread memory.
	caller(func() {
		jqState, err = libjq.NewJqState()
	})
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(jqState).ShouldNot(BeNil())

	// Compile program using created state.
	caller(func() {
		err = jqState.Compile(testProgram)
	})
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(jqState).ShouldNot(BeNil())

	// Process data with compiled program.
	caller(func() {
		defer jqState.Teardown()
		testResult, err = jqState.ProcessOneValue(testData, false)
	})
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(jqState).ShouldNot(BeNil())

	g.Expect(testResult).Should(Equal(testExpected))
}

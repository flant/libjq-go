package jq

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_CgoCall(t *testing.T) {
	g := NewWithT(t)

	in := `{"foo":"baz","bar":"quux"}`

	res, err := NewJq().Program(".").Run(in)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(in))

	g.Expect(cgoCallsCh).ToNot(BeNil(), "cgo calls channel should not be nil after first run")

	res, err = NewJq().Program(".").Run(in)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(in))

	res, err = NewJq().Program(".").Run(in)
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(res).To(Equal(in))
}

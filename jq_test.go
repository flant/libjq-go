package libjq_go

import (
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

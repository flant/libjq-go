package jq

import (
	"sync"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/flant/libjq-go/pkg/libjq"
)

func Test_JqCache_Get_Set(t *testing.T) {
	g := NewWithT(t)

	s := &libjq.JqState{}

	c := NewJqCache()
	c.Set("state", s)

	s2 := c.Get("state")
	g.Expect(s2).Should(Equal(s))
}

func Test_JqCache_Get_Set_Parallel(t *testing.T) {
	g := NewWithT(t)

	s1 := &libjq.JqState{}
	s2 := &libjq.JqState{}
	s3 := &libjq.JqState{}

	c := NewJqCache()
	c.Set("state1", s1)
	c.Set("state2", s2)
	c.Set("state3", s3)

	var wg sync.WaitGroup
	wg.Add(30)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				c.Set("state1", s1)
				c.Set("state2", s2)
				s := c.Get("state3")
				g.Expect(s).Should(Equal(s3))
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				c.Set("state1", s1)
				c.Set("state3", s3)
				s := c.Get("state2")
				g.Expect(s).Should(Equal(s2))
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				c.Set("state3", s3)
				c.Set("state2", s2)
				s := c.Get("state1")
				g.Expect(s).Should(Equal(s1))
			}
		}()
	}

	wg.Wait()

	g.Expect(c.Get("state1")).Should(Equal(s1))
	g.Expect(c.Get("state2")).Should(Equal(s2))
	g.Expect(c.Get("state3")).Should(Equal(s3))
}

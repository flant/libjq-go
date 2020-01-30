package libjq_go

import (
	"github.com/flant/libjq-go/pkg/jq"
)

// Jq is handy shortcut to use a default jq invoker with enabled cache for programs
func Jq() *jq.Jq {
	return jq.NewJq().
		WithCache(jq.JqDefaultCache())
}

package libjq_go

import "github.com/flant/libjq-go/pkg/jq"

// Jq is a default jq invoker with a cache for programs and a jq calls proxy
func Jq() *jq.Jq {
	return jq.NewJq().
		WithCache(jq.JqDefaultCache()).
		WithCallProxy(jq.JqCall)
}

func JqMainThread() *jq.Jq {
	return jq.NewJq().
		WithCache(jq.JqDefaultCache())
}

// JqCallLoop should be called from main.main method to make Jq.Program.Run calls from go routines.
//
// Note: this method locks thread execution.
func JqCallLoop(done chan struct{}) {
	jq.JqCallLoop(done)
}

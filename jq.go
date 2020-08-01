package libjq_go

import (
	"sync"

	"github.com/flant/libjq-go/pkg/jq"
)

var initOnce sync.Once
var defaultCgoCaller jq.CgoCaller
var defaultCache *jq.JqCache

// Jq is a handy shortcut to create a jq invoker with default settings.
//
// Created invokers will share a cache for programs and a cgo caller.
func Jq() *jq.Jq {
	initOnce.Do(func() {
		defaultCache = jq.NewJqCache()
		defaultCgoCaller = jq.NewCgoCaller()
	})
	return jq.NewJq().
		WithCache(defaultCache).
		WithCgoCaller(defaultCgoCaller)
}

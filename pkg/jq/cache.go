package jq

import (
	"sync"

	"github.com/flant/libjq-go/pkg/libjq"
)

/*
Simple cache for jq state objects with compiled programs.
*/

type JqCache struct {
	StateCache map[string]*libjq.JqState
	m          sync.Mutex
}

var jqDefaultCacheInstance *JqCache

var JqDefaultCache = func() *JqCache {
	if jqDefaultCacheInstance == nil {
		jqDefaultCacheInstance = NewJqCache()
	}
	return jqDefaultCacheInstance
}

func NewJqCache() *JqCache {
	return &JqCache{
		StateCache: make(map[string]*libjq.JqState),
		m:          sync.Mutex{},
	}
}

func (jc *JqCache) Get(key string) *libjq.JqState {
	jc.m.Lock()
	defer jc.m.Unlock()
	if v, ok := jc.StateCache[key]; ok {
		return v
	}
	return nil
}

func (jc *JqCache) Set(key string, state *libjq.JqState) {
	jc.m.Lock()
	jc.StateCache[key] = state
	jc.m.Unlock()
}

func (jc *JqCache) Teardown(key string) {
	jc.m.Lock()
	defer jc.m.Unlock()
	if v, ok := jc.StateCache[key]; ok {
		v.Teardown()
	}
}

func (jc *JqCache) TeardownAll() {
	jc.m.Lock()
	defer jc.m.Unlock()
	for _, state := range jc.StateCache {
		state.Teardown()
	}
}

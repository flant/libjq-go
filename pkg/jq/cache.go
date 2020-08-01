package jq

import (
	"sync"

	"github.com/flant/libjq-go/pkg/libjq"
)

// JqCache is a simple cache for JqState objects.
type JqCache struct {
	StateCache map[string]*libjq.JqState
	m          sync.Mutex
}

func NewJqCache() *JqCache {
	return &JqCache{
		StateCache: make(map[string]*libjq.JqState),
		m:          sync.Mutex{},
	}
}

// Get returns cached JqState object or nil of no object is registered for key.
func (jc *JqCache) Get(key string) *libjq.JqState {
	jc.m.Lock()
	defer jc.m.Unlock()
	if v, ok := jc.StateCache[key]; ok {
		return v
	}
	return nil
}

// Set register a JqState object for key.
func (jc *JqCache) Set(key string, state *libjq.JqState) {
	jc.m.Lock()
	jc.StateCache[key] = state
	jc.m.Unlock()
}

// Teardown calls Teardown for cached JqState object.
func (jc *JqCache) Teardown(key string) {
	jc.m.Lock()
	defer jc.m.Unlock()
	if jqState, ok := jc.StateCache[key]; ok {
		jqState.Teardown()
	}
}

// TeardownAll calls Teardown for all cached JqState objects.
func (jc *JqCache) TeardownAll() {
	jc.m.Lock()
	defer jc.m.Unlock()
	for _, state := range jc.StateCache {
		state.Teardown()
	}
}

package libjq

/*
#include <stdlib.h>
#include <jq.h>
*/
import "C"

import (
	"errors"
	"fmt"
)

type JqState struct {
	state *C.struct_jq_state
}

func NewJqState() (*JqState, error) {
	jq := &JqState{}
	jq.state = C.jq_init()

	if jq.state == nil {
		// According to jq sources there is only one error emitter: malloc.
		return nil, fmt.Errorf("cannot allocate new libjq instance")
	}

	return jq, nil
}

func (jq *JqState) SetLibraryPath(path string) {
	if path == "" {
		return
	}
	jvPath := NewJvString(path)
	searchPaths := NewJvArray()
	searchPaths = JvArray(searchPaths.Append(Jv(jvPath)))

	attr := NewJvString("JQ_LIBRARY_PATH")
	C.jq_set_attr(jq.state, attr.v, searchPaths.v)
}

func (jq *JqState) Compile(program string) error {
	result := C.jq_compile(jq.state, C.CString(program))
	if result == 0 {
		return fmt.Errorf("failed to compile: %v", C.jq_get_error_message(jq.state))
	} else {
		return nil
	}
}

func (jq *JqState) CompileArgs(program string, jv Jv) error {
	result := C.jq_compile_args(jq.state, C.CString(program), jv.value())
	if result == 0 {
		return errors.New("failed to compile args")
	} else {
		return nil
	}
}

func (jq *JqState) Start(jv Jv) {
	C.jq_start(jq.state, jv.v, 0)
}

func (jq *JqState) Iterate(f func(Jv)) {
	for {
		value := C.jq_next(jq.state)
		if C.jv_is_valid(value) == 1 {
			f(NewJv(value))
		} else {
			return
		}
	}
}

func (jq *JqState) Teardown() {
	C.jq_teardown(&jq.state)
}

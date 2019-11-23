package libjq

/*
#include <stdlib.h>
#include <jq.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

type JqState struct {
	state *C.struct_jq_state
}

func NewJqState() (*JqState, error) {
	jq := &JqState{}
	jq.state = C.jq_init()

	if jq.state == nil {
		// According to jq sources there is only one error emitter: malloc.
		return nil, fmt.Errorf("cannot allocate memory for the new jq state")
	}

	return jq, nil
}

func (jq *JqState) SetLibraryPath(path string) {
	if path == "" {
		return
	}

	C.jq_set_attr(jq.state, JvString("JQ_LIBRARY_PATH"), JvArray(JvString(path)))
}

func (jq *JqState) Compile(program string) error {
	cProgram := C.CString(program)
	defer C.free(unsafe.Pointer(cProgram))
	result := C.jq_compile(jq.state, cProgram)
	if result == 0 {
		return fmt.Errorf("failed to compile: %v", C.jq_get_error_message(jq.state))
	} else {
		return nil
	}
}

func (jq *JqState) CompileArgs(program string, argsJv C.jv) error {
	cProgram := C.CString(program)
	defer C.free(unsafe.Pointer(cProgram))
	result := C.jq_compile_args(jq.state, cProgram, argsJv)
	if result == 0 {
		return errors.New("failed to compile args")
	} else {
		return nil
	}
}

func (jq *JqState) Teardown() {
	C.jq_teardown(&jq.state)
}

// ProcessOneValue run a jq program over one json object.
//
// Stream mode is not supported: stream splitting can be done by Go code.
func (jq *JqState) ProcessOneValue(inJson string, rawMode bool) (string, error) {
	// Step 1. Validate input data
	cInJson := C.CString(string(inJson))
	defer C.free(unsafe.Pointer(cInJson))
	inputJv := C.jv_parse(cInJson)
	if C.jv_is_valid(inputJv) == 0 {
		return "", JqFormatError(inputJv)
	}

	// defer C.jv_free(inputJv) -> this is not needed:
	//
	// jq_start does two things:
	// - calls jq_reset that free everything on the stack
	// - push inputJv's pointer to the stack
	//
	// If state is used just once, then inputJv will be freed
	// by a Teardown call that will call jq_reset.
	//
	// If state is cached, then inputJv's memory will be freed
	// on the next jq_start call. Also next compile or teardown will
	// lead to jq_reset and memory will be freed.
	// It means that the stack is not freed between the calls, but
	// jq_reset declared as static. So jq_start with jv_null is called
	// to reset the stack between the processing sessions.

	// Step 2. Start processing.
	C.jq_start(jq.state, inputJv, C.int(0))
	defer C.jq_start(jq.state, C.jv_null(), C.int(0))

	// Step 3. Read results.
	out := make([]string, 0)

	var tmp C.jv
	for tmp = C.jq_next(jq.state); C.jv_is_valid(tmp) == 1; tmp = C.jq_next(jq.state) {
		var str string

		if rawMode && C.jv_get_kind(tmp) == C.JV_KIND_STRING {
			str = C.GoString(C.jv_string_value(tmp))
		} else {
			str = JvDumpString(tmp)
		}

		out = append(out, str)

		C.jv_free(tmp)
	}
	return strings.Join(out, "\n"), JqInvalidError(tmp)
}

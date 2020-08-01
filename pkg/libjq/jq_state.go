package libjq

/*
#include <assert.h>
#include <string.h>
#include <stdlib.h>
#include <jv.h>
#include <jq.h>

// jq_format_error without printing to stderr
jv libjq_go_format_error(jv msg) {
  if (jv_get_kind(msg) == JV_KIND_NULL ||
      (jv_get_kind(msg) == JV_KIND_INVALID && !jv_invalid_has_msg(jv_copy(msg)))) {
    jv_free(msg);
    return jv_string("jq: error: out of memory");
  }

  if (jv_get_kind(msg) == JV_KIND_STRING)
    return msg;                         // expected to already be formatted

  if (jv_get_kind(msg) == JV_KIND_INVALID)
    msg = jv_invalid_get_msg(msg);

  if (jv_get_kind(msg) == JV_KIND_NULL)
    return libjq_go_format_error(msg);        // ENOMEM

  // Invalid with msg; prefix with "jq: error: "

  if (jv_get_kind(msg) != JV_KIND_INVALID) {
    if (jv_get_kind(msg) == JV_KIND_STRING)
      return jv_string_fmt("jq: error: %s", jv_string_value(msg));

    msg = jv_dump_string(msg, JV_PRINT_INVALID);
    if (jv_get_kind(msg) == JV_KIND_STRING)
      return jv_string_fmt("jq: error: %s", jv_string_value(msg));
    return libjq_go_format_error(jv_null());  // ENOMEM
  }

  // An invalid inside an invalid!
  return libjq_go_format_error(jv_invalid_get_msg(msg));

}

// Appends error message string to JV_ARRAY passed in data.
static void libjq_go_err_cb(void *data, jv msg) {
  assert(jv_get_kind(*((jv*)data)) == JV_KIND_ARRAY);
  msg = jq_format_error(msg);
  *((jv*)data) = jv_array_append(jv_copy(*((jv*)data)), jv_copy(msg));
  jv_free(msg);
}

// A wrapper around jq_compile to catch compile errors.
// The idea borrowed from jq_test.c
int libjq_go_compile(jq_state *jq, const char* str, jv *msgs) {
    jq_set_error_cb(jq, libjq_go_err_cb, msgs);
	int compiled = jq_compile(jq, str);
	jq_set_error_cb(jq, NULL, NULL);
	return compiled;
}

*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// JqState is a thin wrapper for jq_init, jq_set_attr, jq_compile and jq_start.
//
// It is responsibility of a higher level to call JqState methods in one thread
// as libjq is not compatible with Go's thread migration.
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

// Compile
func (jq *JqState) Compile(program string) error {
	cProgram := C.CString(program)
	defer C.free(unsafe.Pointer(cProgram))

	errMsgs := C.jv_array()
	defer C.jv_free(errMsgs)

	result := C.libjq_go_compile(jq.state, cProgram, &errMsgs)
	if result == 0 {
		msgs := JvArrayToGo(errMsgs)
		return fmt.Errorf("compile: %v", strings.Join(msgs, "\n"))
	} else {
		return nil
	}
}

// TODO create wrapper with errors catching
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

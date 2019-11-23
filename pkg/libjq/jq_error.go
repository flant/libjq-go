package libjq

/*
#include <jq.h>
*/
import "C"
import "errors"

// JqFormatError converts error in jv to Go's error
func JqFormatError(v C.jv) error {
	// jq_format_error prints to stderr "jq: error: out of memory"
	if C.jv_get_kind(v) == C.JV_KIND_NULL {
		return nil
	}

	msg := C.jq_format_error(v)
	defer C.jv_free(msg)

	err := errors.New(C.GoString(C.jv_string_value(msg)))
	return err
}

// JqInvalidError returns error if jv contains invaluid message
func JqInvalidError(jv C.jv) error {
	// jv_invalid_get_msg frees jv.
	msg := C.jv_invalid_get_msg(jv)
	defer C.jv_free(msg)

	switch C.jv_get_kind(msg) {
	case C.JV_KIND_NULL:
		return nil
	case C.JV_KIND_STRING:
		return errors.New(C.GoString(C.jv_string_value(msg)))
	default:
		return errors.New(JvDumpString(msg))
	}
}

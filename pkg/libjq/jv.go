package libjq

/*
#include <stdlib.h>
#include <jv.h>
*/
import "C"

// JvDumpString
func JvDumpString(str C.jv) string {
	// jv_dump_string calls jv_dump_term that frees the provided jv, so we copy it.
	dumpedjv := C.jv_dump_string(C.jv_copy(str), C.int(0))
	defer C.jv_free(dumpedjv)

	// jv_string_value returns a pointer from jv struct,
	// so deferred jv_free for dumpedjv is enough.
	return C.GoString(C.jv_string_value(dumpedjv))
}

// JvString copies Go string to C and return a jv_string
func JvString(str string) C.jv {
	return C.jv_string(C.CString(str))
}

// JvArray returns a jv array value. jq sources has JV_ARRAY macros for this.
func JvArray(first C.jv, items ...C.jv) C.jv {
	arr := C.jv_array()
	arr = C.jv_array_append(arr, first)
	for _, item := range items {
		arr = C.jv_array_append(arr, item)
	}
	return arr
}

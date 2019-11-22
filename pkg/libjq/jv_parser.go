package libjq

/*
#include <jv.h>
*/
import "C"

type JVParser struct {
	V *C.struct_jv_parser
}

func NewJVParser(flags int) JVParser {
	return JVParser{V: C.jv_parser_new(C.int(flags))}
}

func (jvp JVParser) SetBuffer(buffer string) {
	C.jv_parser_set_buf(jvp.V, C.CString(buffer), C.int(len(buffer)), 0)
}

func (jvp JVParser) Iterate(f func(v Jv)) {
	for {
		value := C.jv_parser_next(jvp.V)
		if C.jv_is_valid(value) == 1 {
			f(NewJv(value))
		} else {
			return
		}
	}
}

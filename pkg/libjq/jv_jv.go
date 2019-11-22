package libjq

/*
#include <jv.h>
*/
import "C"

type Jv struct {
	v C.jv
}

func (jv Jv) copy() {
	C.jv_copy(jv.v)
}
func (jv Jv) Free() {
	C.jv_free(jv.v)
}
func (jv Jv) value() C.jv {
	return jv.v
}
func (jv Jv) String() string {
	jv.copy()
	dumped := C.jv_dump_string(jv.v, 0)
	str := C.jv_string_value(dumped)
	//return fmt.Sprintf("%s(%d)", C.GoString(str), jv.RefCount())
	return C.GoString(str)
}

type JvKind C.jv_kind

const (
	KIND_INVALID JvKind = iota
	KIND_NULL
	KIND_FALSE
	KIND_TRUE
	KIND_NUMBER
	KIND_STRING
	KIND_ARRAY
	KIND_OBJECT
)

func (jv Jv) Kind() JvKind {
	return JvKind(C.jv_get_kind(jv.v))
}

func NewJv(v C.jv) Jv {
	return Jv{v: v}
}

func (jv Jv) RefCount() int {
	return int(C.jv_get_refcnt(jv.v))
}

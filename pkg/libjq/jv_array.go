package libjq

/*
#include <jv.h>
*/
import "C"

type JvArray Jv

/// API

func NewJvArray() JvArray {
	return JvArray(NewJv(C.jv_array()))
}

func NewJvArraySized(size int) JvArray {
	return JvArray(NewJv(C.jv_array_sized(C.int(size))))
}

func (jv JvArray) Length() int {
	Jv(jv).copy()
	return int(C.jv_array_length(jv.v))
}

func (jv JvArray) Get(i int) Jv {
	Jv(jv).copy()
	return NewJv(C.jv_array_get(jv.v, C.int(i)))
}

func (jv JvArray) Set(i int, v Jv) Jv {
	Jv(jv).copy()
	return NewJv(C.jv_array_set(jv.v, C.int(i), v.v))
}

func (jv JvArray) Append(v Jv) Jv {
	Jv(jv).copy()
	return NewJv(C.jv_array_append(jv.v, v.v))
}

func (jv JvArray) Concat(v Jv) Jv {
	Jv(jv).copy()
	return NewJv(C.jv_array_concat(jv.v, v.v))
}
func (jv JvArray) Slice(i1 int, i2 int) Jv {
	Jv(jv).copy()
	return NewJv(C.jv_array_slice(jv.v, C.int(i1), C.int(i2)))
}
func (jv JvArray) Indexes(v Jv) Jv {
	Jv(jv).copy()
	return NewJv(C.jv_array_indexes(jv.v, v.v))
}

/// API

/// Helper
func (jv JvArray) Array() (r []Jv) {
	for i := 0; i < jv.Length(); i++ {
		v := jv.Get(i)
		r = append(r, v)
	}
	return
}

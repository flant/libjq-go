package libjq

/*
#include <jv.h>
*/
import "C"

type JvString Jv

func NewJvString(s string) JvString {
	return JvString(NewJv(C.jv_string(C.CString(s))))
}

func NewJvStringSized(s string, size int) JvString {
	return JvString(NewJv(C.jv_string_sized(C.CString(s), C.int(size))))
}

func NewJvStringEmpty(size int) JvString {
	return JvString(NewJv(C.jv_string_empty(C.int(size))))
}

func (jv JvString) LengthBytes() int {
	Jv(jv).copy()
	return int(C.jv_string_length_bytes(jv.v))
}

func (jv JvString) LengthCodePoints() int {
	Jv(jv).copy()
	return int(C.jv_string_length_codepoints(jv.v))
}

func (jv JvString) Hash() uint64 {
	Jv(jv).copy()
	return uint64(C.jv_string_hash(jv.v))
}

func (jv JvString) StringValue() string {
	Jv(jv).copy()
	Jv(jv).copy()
	return C.GoString(C.jv_string_value(jv.v))
}

/*
jv jv_string_indexes(jv j, jv k);
jv jv_string_slice(jv j, int start, int end);
jv jv_string_concat(jv, jv);
jv jv_string_vfmt(const char*, va_list) JV_VPRINTF_LIKE(1);
jv jv_string_fmt(const char*, ...) JV_PRINTF_LIKE(1, 2);
jv jv_string_append_codepoint(jv a, uint32_t c);
jv jv_string_append_buf(jv a, const char* buf, int len);
jv jv_string_append_str(jv a, const char* str);
jv jv_string_split(jv j, jv sep);
jv jv_string_explode(jv j);
jv jv_string_implode(jv j);
*/

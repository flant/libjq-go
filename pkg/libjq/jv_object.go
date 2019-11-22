package libjq

/*#include <jv.h>
 */
import "C"

type JvObject Jv

func NewJvObject() JvObject {
	return JvObject(NewJv(C.jv_object()))
}

func (jv JvObject) Get(key string) Jv {
	Jv(jv).copy()
	panic("impl")
}

func (jv JvObject) Set(key string, value Jv) {
	Jv(jv).copy()
	panic("impl")
}

func (jv JvObject) Delete(key string) {
	Jv(jv).copy()
	panic("impl")
}

func (jv JvObject) Length() int {
	Jv(jv).copy()
	return int(C.jv_object_length(jv.v))
}

func (jv JvObject) Merge(target JvObject) JvObject {
	Jv(jv).copy()
	return JvObject(NewJv(C.jv_object_merge(jv.v, target.v)))
}

func (jv JvObject) MergeRecursive(target JvObject) JvObject {
	Jv(jv).copy()
	return JvObject(NewJv(C.jv_object_merge_recursive(jv.v, target.v)))
}

func (jv JvObject) Iter() int {
	return int(C.jv_object_iter(jv.v))
}

func (jv JvObject) IterNext(i int) int {
	return int(C.jv_object_iter_next(jv.v, C.int(i)))
}

func (jv JvObject) IterValid(i int) bool {
	return int(C.jv_object_iter_valid(jv.v, C.int(i))) == 1
}

func (jv JvObject) IterKey(i int) Jv {
	return NewJv(C.jv_object_iter_key(jv.v, C.int(i)))
}

func (jv JvObject) IterValue(i int) Jv {
	return NewJv(C.jv_object_iter_value(jv.v, C.int(i)))
}

/*
jv jv_object_get(jv object, jv key);
jv jv_object_set(jv object, jv key, jv value);
jv jv_object_delete(jv object, jv key);

*/

func (jv JvObject) ForEach(c func(Jv, Jv)) {
	iter := jv.Iter()
	for jv.IterValid(iter) {
		k := jv.IterKey(iter)
		v := jv.IterValue(iter)
		c(k, v)
		iter = jv.IterNext(iter)
	}
}

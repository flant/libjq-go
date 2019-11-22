package libjq_go

import "github.com/flant/libjq-go/jq"

func Jq() *jq.Jq {
	return jq.NewJq()
}

func JqWithLib(path string) *jq.Jq {
	njq :=  jq.NewJq()
	njq.LibPath = path
	return njq
}

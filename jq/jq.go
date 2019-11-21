package jq

type Jq struct {
	LibPath string
}

func NewJq() *Jq {
	return &Jq{}
}


func (jq *Jq) Program(program string) *JqProgram {
	return &JqProgram{
		Jq: jq,
		Program: program,
	}
}

type JqProgram struct {
	Jq *Jq
	Program string
}


func (jqp *JqProgram) Run(data string) (string, error) {
	// lib := jqp.Jq.LibPath
	return "", nil
}

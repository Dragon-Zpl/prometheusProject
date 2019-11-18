package prometheus

import "sync"

type RegisterProm struct {
	Vec *Prom
	JobName string
	Lock sync.Mutex
}

var RegisterPromMap map[string]*RegisterProm

func init() {
	RegisterPromMap = make(map[string]*RegisterProm)
}

func NewRegisterProm() *RegisterProm {

	return &RegisterProm{}
}
// SPDX-License-Identifier: MIT

package pipe

import "sync"

type Func = func()

type Dispatcher struct {
	m     sync.Mutex
	funcs []Func
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		funcs: make([]Func, 0, 10),
	}
}

func (d *Dispatcher) Add(f Func) {
	d.m.Lock()
	defer d.m.Unlock()
	d.funcs = append(d.funcs, f)
}

func (d *Dispatcher) Run() {
	d.m.Lock()
	defer d.m.Unlock()

	for _, f := range d.funcs {
		f()
	}

	d.funcs = d.funcs[:0]
}

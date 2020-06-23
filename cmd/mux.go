package cmd

import (
	"github.com/soypat/decimate/csvtools"
)

type Mux struct {
	In  *chan csvtools.Value
	Out []*chan csvtools.Value
}

func (m *Mux) run() {
	var EOF,ok bool
	var val csvtools.Value
	for !EOF {
		val, ok = <- *m.In
		if ok {
			for i, _ := range m.Out {
				*m.Out[i] <- val
			}
			if val.Type() == csvtools.TypeEOF {
				EOF = true
			}
		}
	}
}

func (m *Mux) makeOutputs(quantity int) {
	m.Out = make([]*chan csvtools.Value, quantity)
	for i, _ := range m.Out {
		c := make(chan csvtools.Value, bufferSize)
		m.Out[i] = &c
	}
}

func (m *Mux) setOutputs(cs []chan csvtools.Value) {
	m.Out = make([]*chan csvtools.Value, len(cs))
	for i, _ := range cs {
		m.Out[i] = &cs[i]
	}
}

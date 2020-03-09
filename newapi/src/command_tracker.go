package main

import (
	"bytes"
)

type trackable struct {
	length int
	chain  []byte
	run    func()
}

type Tracker struct {
	currentType []byte
	trackables  []trackable
}

func (t *Tracker) Add(chain []byte, run func()) {
	var newtrack trackable

	newtrack.chain = chain
	newtrack.length = len(chain)
	newtrack.run = run
	t.trackables = append(t.trackables, newtrack)
}

func (t *Tracker) Push(c byte) {
	t.currentType = append(t.currentType, c)
	typeLen := len(t.currentType)
	maxTrLength := 0
	for _, tr := range t.trackables {
		from := typeLen - tr.length
		if maxTrLength < tr.length {
			maxTrLength = tr.length
		}
		if from < 1 {
			continue
		}
		if bytes.Compare(t.currentType[from:], tr.chain) == 0 {
			tr.run()
		}
	}
	from := typeLen - maxTrLength
	if from > 0 {
		t.currentType = t.currentType[from:]
	}
}

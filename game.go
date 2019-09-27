package main

import (
	"bytes"

	svg "github.com/ajstarks/svgo/float"
)

var turns = make(map[string]Turn)

// Turn stores the SVG data as a buffer ready to be sent in a web request
type Turn struct {
	bytes.Buffer
}

func NewTurn() *Turn {
	return &Turn{
		Buffer: bytes.Buffer{},
	}
}

func (t Turn) Svg() []byte {
	buf := &bytes.Buffer{}
	canvas := svg.New(buf)
	canvas.Startraw("class=\"fillscreen\"")
	canvas.Gid("scene")
	buf.Write(t.Bytes())
	canvas.Gend()
	canvas.End()
	return buf.Bytes()
}

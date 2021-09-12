package io

import (
	"bytes"
	. "context"
)

type Writer interface {
	Write(b []byte) (int, error)
}

func Get() Writer {
	var b bytes.Buffer
	return &b
}

func Context() Context {
	return Background()
}

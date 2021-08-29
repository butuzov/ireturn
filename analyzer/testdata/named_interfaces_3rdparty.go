package example

import "github.com/foo/bar"

func s() bar.Buzzer {
	return nil
}

type f struct{}

func (frec f) Buzz() {}

func (frec f) New() bar.Buzzer {
	return &f{}
}

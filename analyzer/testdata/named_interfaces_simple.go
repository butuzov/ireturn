package example

type (
	iDoer         interface{}
	iDoerAny      any
	iDoerConcreat int
)

func newIDoer() iDoer {
	return iDoerConcreat(0)
}

func newIDoerAny() iDoerAny {
	return iDoerConcreat(0)
}

type Fooer interface {
	Foo()
}

type Barer interface {
	Bar()
}

type FooerBarer interface {
	Fooer
	Barer
}

type nameStruct struct{}

func NewNamedStruct() FooerBarer {
	return &nameStruct{}
}

func (ns nameStruct) Foo() {}
func (ns nameStruct) Bar() {}

package example

type _t_anon_interface int

func (t *_t_anon_interface) Do() {}

func NewAnonymousInterface() interface {
	Do()
} {
	t := _t_anon_interface(1)
	return &t
}

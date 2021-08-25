package ireturn

type itype int

const (
	typeEmptyInterface itype = 1 << iota
	typeAnonInterface
	typeErrorInterface
	typeNamedInterface
)

type iface struct {
	name string // preserved for named interfaces
	pos  int    // position in return tuple
	t    itype  // type of the interface
}

func issue(name string, pos int, interfaceType itype) iface {
	return iface{
		name: name,
		pos:  pos,
		t:    interfaceType,
	}
}

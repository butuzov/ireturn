package ireturn

type itype uint8

const (
	typeEmptyInterface itype = 1 << iota // ref as empty
	typeAnonInterface                    // ref as anon
	typeErrorInterface                   // ref as error
	typeNamedInterface                   // ref as named
)

const (
	nameEmpty = "empty"
	nameAnon  = "anon"
	nameError = "error"
	// TODO(butuzov): pic stdlib interfaces
	// nameStdLib = "stdlib"
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

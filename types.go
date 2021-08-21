package ireturn

type itype int

const (
	typeEmptyInterface itype = 1 << iota
	typeAnonInterface
)

type iface struct {
	name string // preserved for named iterfaces
	pos  int    // position in return tuple
	t    itype  // type of the interface
}

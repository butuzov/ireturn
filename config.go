package ireturn

// you can use either allow or reject
// you can allow standard library + some extra stuff to be used (error / interface)

// Config allows linter to be configurable.

type Action int

const (
	Allow Action = iota
	Reject
)

type Config struct {
	Action Action
	List   []string
}

func NewDefaultConfig() Config {
	return Config{
		Action: Allow,
		List:   []string{},
	}
}

// https://github.com/tomarrell/wrapcheck/blob/master/wrapcheck/wrapcheck.go
// https://github.com/ldez/tagliatelle

package ireturn

// you can use either allow or reject
// you can allow standard library + some extra stuff to be used (error / interface)

// Config allows linter to be configurable.

type Action int

const (
	Allow Action = iota // default action for empty zero config... reject none.
	Reject
)

type Config struct {
	Action Action
	List   []string

	// private fields (for search optimization look ups)
	init  bool
	quick uint8
}

func NewDefaultConfig() Config {
	return Config{
		Action: Allow,
		List: []string{
			nameEmpty,
			nameError,
			nameAnon,
		},
	}
}

func (config *Config) Has(i iface) bool {
	if !config.init {
		config.compileList()
		config.init = true
	}

	if config.quick&uint8(i.t) > 0 {
		return true
	}
	return false
}

// compileList will transform text list into a bitmask for quick searches and
// slice of regular expressions for quick searches.
func (config *Config) compileList() {
	for _, str := range config.List {
		switch str {
		case nameError:
			config.quick |= uint8(typeErrorInterface)
		case nameEmpty:
			config.quick |= uint8(typeEmptyInterface)
		case nameAnon:
			config.quick |= uint8(typeAnonInterface)
		}
		// todo(butuzov): named interfaces check?

		// todo(butuzov): how can we log error in golangci-lint?
	}
}

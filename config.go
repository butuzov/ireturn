package ireturn

import (
	"regexp"
)

// defaultConfig is core of the validation, ...
// todo(butuzov): write proper intro...

type defaultConfig struct {
	List []string

	// private fields (for search optimization look ups)
	init  bool
	quick uint8
	list  []*regexp.Regexp
}

func (config *defaultConfig) Has(i iface) bool {
	if !config.init {
		config.compileList()
		config.init = true
	}

	if config.quick&uint8(i.t) > 0 {
		return true
	}

	// not a named interface (because error, interface{}, anon interface has keywords.)
	if i.t&typeNamedInterface == 0 && i.t&typeNamedStdInterface == 0 {
		return false
	}

	for _, re := range config.list {
		if re.MatchString(i.name) {
			return true
		}
	}

	return false
}

// compileList will transform text list into a bitmask for quick searches and
// slice of regular expressions for quick searches.
func (config *defaultConfig) compileList() {
	for _, str := range config.List {
		switch str {
		case nameError:
			config.quick |= uint8(typeErrorInterface)
		case nameEmpty:
			config.quick |= uint8(typeEmptyInterface)
		case nameAnon:
			config.quick |= uint8(typeAnonInterface)
		case nameStdLib:
			config.quick |= uint8(typeNamedStdInterface)
		}

		// allow to parse regular expressions
		// todo(butuzov): how can we log error in golangci-lint?
		if re, err := regexp.Compile(str); err == nil {
			config.list = append(config.list, re)
		}

	}
}

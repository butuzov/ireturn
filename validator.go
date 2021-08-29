package ireturn

type validator interface {
	isValid(iface) bool
}

// rejectConfig specifies a list of interfaces (keywords, patters and regular expressions)
// that are rejected by ireturn as valid to return, any non listed interface are allowed.
type rejectConfig struct {
	*defaultConfig
}

func RejectAll(patterns []string) *rejectConfig {
	return &rejectConfig{&defaultConfig{List: patterns}}
}

func (rc *rejectConfig) isValid(i iface) bool {
	return !rc.Has(i)
}

// allowConfig specifies a list of interfaces (keywords, patters and regular expressions)
// that are allowed by ireturn as valid to return, any non listed interface are rejected.
type allowConfig struct {
	*defaultConfig
}

func AllowAll(patterns []string) *allowConfig {
	return &allowConfig{&defaultConfig{List: patterns}}
}

func (ac *allowConfig) isValid(i iface) bool {
	return ac.Has(i)
}

//nolint: exhaustivestruct
func DefaultValidatorConfig() *allowConfig {
	return AllowAll([]string{
		nameEmpty,  // "empty" - for all empty interfaces (interface{})
		nameError,  // "error" - for all error's
		nameAnon,   // "anon" - for all empty interfaces with methods (interface {Method()})
		nameStdLib, // "std" - for all standard library packages
	})
}

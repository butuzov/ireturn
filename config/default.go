package config

import "github.com/butuzov/ireturn/types"

//nolint: exhaustivestruct
func DefaultValidatorConfig() *allowConfig {
	return AllowAll([]string{
		types.NameEmpty,  // "empty": empty interfaces (interface{})
		types.NameError,  // "error": for all error's
		types.NameAnon,   // "anon": for all empty interfaces with methods (interface {Method()})
		types.NameStdLib, // "std": for all standard library packages
	})
}

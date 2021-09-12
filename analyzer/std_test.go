package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isStdLib(t *testing.T) {
	tests := map[string]bool{
		"context.Context":                     true,
		"io/fs.File":                          true,
		"github.com/user/pkg/context.Context": false,
		"foo/bar.Context":                     false,
	}

	for name, want := range tests {
		want, name := want, name
		t.Run(name, func(t *testing.T) {
			got := isStdPkgInterface(name)
			assert.Equal(t, got, want,
				"pkg %s doens't match expectations (got %v vs want %v)", name, got, want)
		})
	}
}

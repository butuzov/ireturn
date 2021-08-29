package analyzer

import "testing"

func Test_isStdLib(t *testing.T) {
	tests := map[string]struct {
		name string
		want bool
	}{
		"context.Context": {
			want: true,
		},

		"io/fs.File": {
			want: true,
		},

		"github.com/user/pkg/context.Context": {
			want: false,
		},
	}

	for name, tt := range tests {
		tt := tt
		name := name
		t.Run(name, func(t *testing.T) {
			got := isStdLib(name)
			assert(t, tt.want == got,
				"pkg %s doens't match expectations (got %v vs want %v)", name, got, tt.want)
		})
	}
}

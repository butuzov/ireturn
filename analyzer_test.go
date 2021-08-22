package ireturn

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	tests := []testCase{}
	//
	tests = append(tests, testCase{
		name: "Empty Package With No Issues",
		mask: "hello-world.go",
	})
	//
	tests = append(tests, testCase{
		name: "Empty Interface Return",
		mask: "empty_interface.go",
		want: []string{
			"fooInterface returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "Anonymouse Interface",
		mask: "anonymouse_interafce.go",
		want: []string{
			"NewAnonymouseInterface returns interface (anonymouse interface)",
		},
	})

	tests = append(tests, testCase{
		name: "Correct Disallow Directive",
		mask: "disallow_directive_ok.go",
		want: []string{
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "Error Interface return",
		mask: "errors.go",
		want: []string{
			"errorReturn returns interface (error)",
			"errorAliasReturn returns interface (error)",
			"errorTypeReturn returns interface (error)",
			"newErrorInterface returns interface (error)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interface",
		mask: "named_*.go",
		want: []string{
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, tt.test())
	}
}

// ---------------------------------------------------------- fake test --------

type fakeTest struct{}

func (t *fakeTest) Errorf(format string, args ...interface{}) {}

// ---------------------------------------------------------- test case --------
type testCase struct {
	name string
	mask string // file mask
	want []string
}

func (tc testCase) test() func(*testing.T) {
	return func(t *testing.T) {
		// -------------------------------------------------------------- setup ----
		goroot, dir, err := directory(t)
		if err != nil {
			t.Error(err)
		}

		if err := cp("go.mod", dir); err != nil {
			t.Error(err)
		}

		if err := tc.xerox(dir); err != nil {
			t.Error(err)
		}

		// --------------------------------------------------------------- test ----
		results := analysistest.Run(&fakeTest{}, goroot, NewAnalyzer(), "example")

		// ------------------------------------------------------------ results ----

		var tmp []string
		for _, d := range results[0].Diagnostics {
			tmp = append(tmp, d.Message)
		}

		if diff := cmp.Diff(tc.want, tmp); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func directory(t *testing.T) (goroot, dir string, err error) {
	t.Helper()

	goroot = t.TempDir()
	dir = filepath.Join(goroot, "src/example")

	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", "", err
	}

	return goroot, dir, nil
}

func (tc testCase) xerox(dest string) error {
	matches, err := filepath.Glob("testdata/" + tc.mask)
	if err != nil {
		return err
	}

	for _, file := range matches {
		if err := cp(file, dest); err != nil {
			return err
		}
	}

	return nil
}

func cp(src, dst string) error {
	if location, err := filepath.Abs(src); err != nil {
		return err
	} else if data, err := ioutil.ReadFile(location); err != nil {
		return err
	} else if err := ioutil.WriteFile(filepath.Join(dst, filepath.Base(src)), data, 0600); err != nil {
		return err
	}

	return nil
}

//nolint: unused, deadcode //
func assert(t *testing.T, condHappend bool, msg string, args ...interface{}) bool {
	t.Helper()
	if condHappend {
		return true
	}

	t.Errorf(msg, args...)
	return false
}

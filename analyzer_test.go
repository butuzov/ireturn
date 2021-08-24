package ireturn

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/analysis/analysistest"
)

const testPackageName = "example"

func TestAll(t *testing.T) {
	tests := []testCase{}
	//
	tests = append(tests, testCase{
		name: "Empty Package With No Issues",
		mask: []string{
			"hello-world.go",
			"go.*",
		},
	})
	//
	tests = append(tests, testCase{
		name: "Empty Interface Return",
		mask: []string{"empty_interface.go", "go.*"},
		want: []string{
			"fooInterface returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "Anonymouse Interface",
		mask: []string{"anonymouse_interafce.go", "go.*"},
		want: []string{
			"NewAnonymouseInterface returns interface (anonymouse interface)",
		},
	})

	tests = append(tests, testCase{
		name: "Correct Disallow Directive",
		mask: []string{"disallow_directive_ok.go", "go.*"},
		want: []string{
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "Error Interface return",
		mask: []string{"errors.go", "go.*"},
		want: []string{
			"errorReturn returns interface (error)",
			"errorAliasReturn returns interface (error)",
			"errorTypeReturn returns interface (error)",
			"newErrorInterface returns interface (error)",
		},
	})

	// because of https://github.com/golang/go/issues/37054
	// we not going to import external modules in tests,
	// but rather create new ones that are "external"
	tests = append(tests, testCase{
		name: "Named Interface",
		mask: []string{"named_*.go", "github.com/foo/bar/*"},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
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
	mask []string // file mask
	want []string
}

func (tc testCase) test() func(*testing.T) {
	return func(t *testing.T) {
		// -------------------------------------------------------------- setup ----
		goroot, srcdir, err := directory(t)
		if err != nil {
			t.Error(err)
		}

		if err := tc.xerox(srcdir); err != nil {
			t.Error(err)
		}

		// --------------------------------------------------------------- test ----
		results := analysistest.Run(&fakeTest{}, goroot, NewAnalyzer(), testPackageName)

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

func directory(t *testing.T) (goroot, srcdir string, err error) {
	t.Helper()

	goroot = t.TempDir()
	srcdir = filepath.Join(goroot, "src")

	if err := os.MkdirAll(srcdir, 0777); err != nil {
		return "", "", err
	}

	return goroot, srcdir, nil
}

func (tc testCase) xerox(root string) error {
	for _, mask := range tc.mask {

		files, err := filepath.Glob("testdata/" + mask)
		if err != nil {
			return err
		}

		for _, file := range files {
			// directory
			isInSubDir := strings.Count(file, "/") > 1
			directory := testPackageName
			if isInSubDir {
				// cut off suffix & prefix
				directory = file[len("testdata")+1 : len(file)-len(filepath.Base(file))-1]
			}

			// create if no exists
			if err := os.MkdirAll(filepath.Join(root, directory), 0777); err != nil {
				return err
			}

			// copy
			if err := cp(file, filepath.Join(root, directory)); err != nil {
				return err
			}
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

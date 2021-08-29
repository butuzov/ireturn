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

	tests = append(tests, testCase{
		name: "Hello World",
		mask: []string{
			"hello-world.go",
			"go.*",
		},
		want: []string{},
	})

	tests = append(tests, testCase{
		name: "interface{}/allow",
		mask: []string{"empty_interface.go", "go.*"},
		want: []string{},
		reject: AllowAll([]string{
			nameEmpty, // empty
		}),
	})

	tests = append(tests, testCase{
		name: "interface{}/reject",
		mask: []string{"empty_interface.go", "go.*"},
		want: []string{
			"fooInterface returns interface (interface{})",
		},
		reject: RejectAll([]string{
			nameEmpty,
		}),
	})

	tests = append(tests, testCase{
		name:   "Anonymouse Interface/allow",
		mask:   []string{"anonymouse_interafce.go", "go.*"},
		want:   []string{}, // no errors expected as anon interfaces are allowed
		reject: AllowAll([]string{nameAnon}),
	})

	tests = append(tests, testCase{
		name: "Anonymouse Interface/reject",
		mask: []string{"anonymouse_interafce.go", "go.*"},
		want: []string{
			"NewAnonymouseInterface returns interface (anonymouse interface)",
		},
		reject: RejectAll(
			[]string{nameAnon},
		),
	})

	tests = append(tests, testCase{
		name: "Disallow Directives",
		mask: []string{"disallow_directive_ok.go", "go.*"},
		want: []string{
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
		},
		reject: RejectAll([]string{
			nameEmpty,
		}),
	})

	tests = append(tests, testCase{
		name: "Error/allow",
		mask: []string{"errors.go", "go.*"},
		reject: AllowAll([]string{
			nameError, //
		}),
		want: []string{},
	})

	tests = append(tests, testCase{
		name: "Error/reject",
		mask: []string{"errors.go", "go.*"},
		reject: RejectAll([]string{
			nameError,
		}),
		want: []string{
			"errorReturn returns interface (error)",
			"errorAliasReturn returns interface (error)",
			"errorTypeReturn returns interface (error)",
			"newErrorInterface returns interface (error)",
		},
	})

	// 1) because of https://github.com/golang/go/issues/37054
	// we not going (we can't) import external modules in our tests,
	// but rather we will create new ones that are "external"
	// 2) * we can't (and shouldn't) specify named global pattern for named.
	tests = append(tests, testCase{
		name:   "Named Interfaces/(allow*)",
		mask:   []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		reject: AllowAll([]string{}),
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/stdlib/reject",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"NamedContext returns interface (context.Context)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
		reject: RejectAll([]string{
			nameStdLib,
		}),
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/stdlib/allow",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
		reject: AllowAll([]string{
			nameStdLib, //
		}),
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/pattern/allow",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		reject: AllowAll([]string{
			"github.com/foo/bar", // only valid interface is from this package.
		}),
		want: []string{
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/pattern/reject",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		reject: RejectAll([]string{
			"github.com/foo/bar",
		}),
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/regexp/reject",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"NewDeclared returns interface (internal/sample.Doer)",
		},
		reject: RejectAll([]string{
			"\\.Doer", // .Doer only is invalid interface
		}),
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/regexp/allow",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
		reject: AllowAll([]string{
			"\\.Doer", // .Doer only is valid interface
		}),
	})

	tests = append(tests, testCase{
		name: "default/all/reject_all_but_named(non_std)",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
		reject: DefaultValidatorConfig(),
	})

	tests = append(tests, testCase{
		name: "default/all/nil_reject",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
	})

	tests = append(tests, testCase{
		name: "all/stdlib/allow",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		want: []string{
			"NewAnonymouseInterface returns interface (anonymouse interface)",
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
			"fooInterface returns interface (interface{})",
			"errorReturn returns interface (error)",
			"errorAliasReturn returns interface (error)",
			"errorTypeReturn returns interface (error)",
			"newErrorInterface returns interface (error)",
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
		reject: AllowAll([]string{
			nameStdLib,
		}),
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
	name   string
	reject validator
	mask   []string // file mask
	want   []string
}

func (tc testCase) test() func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		// -------------------------------------------------------------- setup ----
		goroot, srcdir, err := directory(t)
		if err != nil {
			t.Error(err)
		}

		if err := tc.xerox(srcdir); err != nil {
			t.Error(err)
		}

		// --------------------------------------------------------------- test ----
		results := analysistest.Run(
			&fakeTest{}, goroot, NewAnalyzerWithConfig(tc.reject), testPackageName)

		// ------------------------------------------------------------ results ----

		tmp := []string{}
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

//nolint: unused
func assert(t *testing.T, condHappend bool, msg string, args ...interface{}) bool {
	t.Helper()
	if condHappend {
		return true
	}

	t.Errorf(msg, args...)
	return false
}

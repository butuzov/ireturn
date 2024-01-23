//nolint: wrapcheck
//nolint: exhaustivestruct

package analyzer

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/butuzov/ireturn/analyzer/internal/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

const testPackageName = "example"

// ---------------------------------------------------------- test case --------

type testCase struct {
	name string            // test name
	meta map[string]string // options
	mask []string          // glob of files/directories we going to analyze
	pkgm string            // package name, if empty defaults to "example"

	// expectations -----
	want []string // list of errors
	fail error    // if linter expected to fail - errors is expected to be non nil
}

func TestAll(t *testing.T) {
	tests := []testCase{}

	tests = append(tests, testCase{
		name: "Hello World",
		mask: []string{"hello-world.go", "go.*"},
		meta: map[string]string{},
		want: []string{},
	})

	// tests = append(tests, testCase{
	// 	name: "Generic Interface",
	// 	mask: []string{"generic.go", "go.*"},
	// 	meta: map[string]string{},
	// 	want: []string{},
	// })

	tests = append(tests, testCase{
		name: "interface{}/allow",
		mask: []string{"empty_interface.go", "go.*"},
		meta: map[string]string{
			"allow": types.NameEmpty,
		},
		want: []string{},
	})

	tests = append(tests, testCase{
		name: "interface{}/reject",
		mask: []string{"empty_interface.go", "go.*"},
		meta: map[string]string{
			"reject": types.NameEmpty,
		},
		want: []string{
			"fooInterface returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "anonymous Interface/allow",
		mask: []string{"anonymous_interafce.go", "go.*"},
		meta: map[string]string{
			"allow": types.NameAnon,
		},
		want: []string{}, // no errors expected as anon interfaces are allowed
	})

	tests = append(tests, testCase{
		name: "anonymous Interface/reject",
		mask: []string{"anonymous_interafce.go", "go.*"},
		meta: map[string]string{
			"reject": types.NameAnon,
		},
		want: []string{
			"NewanonymousInterface returns interface (anonymous interface)",
		},
	})

	tests = append(tests, testCase{
		name: "Disallow Directives",
		mask: []string{"disallow_directive_ok.go", "go.*"},
		meta: map[string]string{
			"reject": types.NameEmpty,
		},
		want: []string{
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "Disallow Directives 2",
		mask: []string{"disallow_directive_ok.go", "go.*"},
		meta: map[string]string{
			"reject":   types.NameEmpty,
			"nonolint": "true",
		},
		want: []string{
			"dissAllowDirective1 returns interface (interface{})",
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective3 returns interface (interface{})",
			"dissAllowDirective4 returns interface (interface{})",
			"dissAllowDirective5 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
		},
	})

	tests = append(tests, testCase{
		name: "Error/allow",
		mask: []string{"errors.go", "go.*"},
		meta: map[string]string{
			"allow": types.NameError,
		},
		want: []string{},
	})

	tests = append(tests, testCase{
		name: "Error/reject",
		mask: []string{"errors.go", "go.*"},
		meta: map[string]string{
			"reject": types.NameError,
		},
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
		name: "Named Interfaces/allow*",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"allow": "",
		},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/stdlib/reject",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"reject": types.NameStdLib,
		},
		want: []string{
			"NamedContext returns interface (context.Context)",
			"NamedBytes returns interface (io.Writer)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/stdlib/allow",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"allow": types.NameStdLib,
		},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/pattern/allow",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"allow": "github.com/foo/bar", // only valid interface is from this package.
		},
		want: []string{
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
			"NamedBytes returns interface (io.Writer)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/pattern/reject",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"reject": "github.com/foo/bar", // only valid interface is from this package.
		},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/regexp/reject",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"reject": "\\.Doer", //
		},
		want: []string{
			"NewDeclared returns interface (internal/sample.Doer)",
		},
	})

	tests = append(tests, testCase{
		name: "Named Interfaces/regexp/allow",
		mask: []string{"named_*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"allow": "\\.Doer", // allow only Doer interfaces from any package
		},
		want: []string{
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
			"NamedBytes returns interface (io.Writer)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
	})

	tests = append(tests, testCase{
		name: "defaults",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{}, // skipping any configuration to run default one.
		want: []string{
			"Min returns generic interface (T) of type param ~int | ~float64 | ~float32",
			"MixedReturnParameters returns generic interface (T) of type param ~int | ~float64 | ~float32",
			"MixedReturnParameters returns generic interface (K) of type param ~int | ~float64 | ~float32",
			"Max returns generic interface (foobar) of type param ~int | ~float64 | ~float32",
			"SumIntsOrFloats returns generic interface (V) of type param int64 | float64",
			"FuncWithGenericAny_NamedReturn returns generic interface (T_ANY) of type param any",
			"FuncWithGenericAny returns generic interface (T_ANY) of type param any",
			"Get returns generic interface (V_COMPARABLE) of type param comparable",
			"Get returns generic interface (V_ANY) of type param any",
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
	})

	tests = append(tests, testCase{
		name: "all/stdlib/allow",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"allow": types.NameStdLib, // allow only interfaces from standard library (e.g. io.Writer, fmt.Stringer)
		},
		want: []string{
			"NewanonymousInterface returns interface (anonymous interface)",
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
			"fooInterface returns interface (interface{})",
			"errorReturn returns interface (error)",
			"errorAliasReturn returns interface (error)",
			"errorTypeReturn returns interface (error)",
			"newErrorInterface returns interface (error)",
			"Min returns generic interface (T) of type param ~int | ~float64 | ~float32",
			"MixedReturnParameters returns generic interface (T) of type param ~int | ~float64 | ~float32",
			"MixedReturnParameters returns generic interface (K) of type param ~int | ~float64 | ~float32",
			"Max returns generic interface (foobar) of type param ~int | ~float64 | ~float32",
			"SumIntsOrFloats returns generic interface (V) of type param int64 | float64",
			"FuncWithGenericAny_NamedReturn returns generic interface (T_ANY) of type param any",
			"FuncWithGenericAny returns generic interface (T_ANY) of type param any",
			"Get returns generic interface (V_COMPARABLE) of type param comparable",
			"Get returns generic interface (V_ANY) of type param any",
			"FunctionAny returns interface (any)",
			"FunctionInterface returns interface (interface{})",
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
		},
	})

	// Rejecting Only Generic Returns
	tests = append(tests, testCase{
		name: "generic/reject",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"reject": types.NameGeneric, // reject only generic interfaces
		},
		want: []string{
			"Min returns generic interface (T) of type param ~int | ~float64 | ~float32",
			"MixedReturnParameters returns generic interface (T) of type param ~int | ~float64 | ~float32",
			"MixedReturnParameters returns generic interface (K) of type param ~int | ~float64 | ~float32",
			"Max returns generic interface (foobar) of type param ~int | ~float64 | ~float32",
			"SumIntsOrFloats returns generic interface (V) of type param int64 | float64",
			"FuncWithGenericAny_NamedReturn returns generic interface (T_ANY) of type param any",
			"FuncWithGenericAny returns generic interface (T_ANY) of type param any",
			"Get returns generic interface (V_COMPARABLE) of type param comparable",
			"Get returns generic interface (V_ANY) of type param any",
		},
	})

	tests = append(tests, testCase{
		name: "generic/allow",
		mask: []string{"*.go", "github.com/foo/bar/*", "internal/sample/*"},
		meta: map[string]string{
			"allow": types.NameGeneric, // allow only generic interfaces
		},
		want: []string{
			"NewanonymousInterface returns interface (anonymous interface)",
			"dissAllowDirective2 returns interface (interface{})",
			"dissAllowDirective6 returns interface (interface{})",
			"fooInterface returns interface (interface{})",
			"errorReturn returns interface (error)",
			"errorAliasReturn returns interface (error)",
			"errorTypeReturn returns interface (error)",
			"newErrorInterface returns interface (error)",
			"FunctionAny returns interface (any)",
			"FunctionInterface returns interface (interface{})",
			"s returns interface (github.com/foo/bar.Buzzer)",
			"New returns interface (github.com/foo/bar.Buzzer)",
			"NewDeclared returns interface (internal/sample.Doer)",
			"newIDoer returns interface (example.iDoer)",
			"newIDoerAny returns interface (example.iDoerAny)",
			"NewNamedStruct returns interface (example.FooerBarer)",
			"NamedContext returns interface (context.Context)",
			"NamedBytes returns interface (io.Writer)",
			"NamedStdFile returns interface (go/types.Importer)",
		},
	})

	tests = append(tests, testCase{
		name: "allow/reject",
		mask: []string{"hello-world.go", "go.*"},
		meta: map[string]string{
			"allow":  types.NameStdLib,
			"reject": types.NameStdLib,
		},
		want: []string{},
	})

	for _, tt := range tests {
		t.Run(tt.name, tt.test())
	}
}

// ---------------------------------------------------------- spy test ---------

type spyTest struct {
	errors []error
}

func (st *spyTest) Errorf(format string, args ...interface{}) {
	st.errors = append(st.errors, fmt.Errorf(format, args...))
}

func (tc testCase) test() func(*testing.T) {
	return func(t *testing.T) {
		// ---------------------------------------------------------- setup ----
		goroot, srcdir, err := directory(t)
		if err != nil {
			t.Error(err)
		}

		if err := tc.xerox(srcdir); err != nil {
			t.Error(err)
		}

		// ----------------------------------------------------------- test ----
		analyzer := NewAnalyzer()

		if len(tc.meta) > 0 {
			fs := flag.NewFlagSet("", flag.ExitOnError)
			for key, value := range tc.meta {
				fs.String(key, value, "")
			}
			analyzer.Flags = *fs
		}

		st := &spyTest{errors: []error{}}
		pkgName := testPackageName
		if tc.pkgm != "" {
			pkgName = tc.pkgm
		}
		results := analysistest.Run(st, goroot, analyzer, pkgName)

		// -------------------------------------------------------- results ----

		tmp := []string{}
		for _, d := range results[0].Diagnostics {
			tmp = append(tmp, d.Message)
		}

		assert.Equal(t, tc.want, tmp)

		// --------------------------------------------------------- errors ----
		if tc.fail != nil {
			for _, err := range st.errors {
				got := err.Error()
				fmt.Println(">", got)
				assert.Containsf(t, got, tc.fail.Error(), "unexpected error: %#v", err)
			}
		}
	}
}

func directory(t *testing.T) (goroot, srcdir string, err error) {
	t.Helper()

	goroot = t.TempDir()
	srcdir = filepath.Join(goroot, "src")

	if err := os.MkdirAll(srcdir, 0o777); err != nil {
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
			if err := os.MkdirAll(filepath.Join(root, directory), 0o777); err != nil {
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
	} else if err := ioutil.WriteFile(filepath.Join(dst, filepath.Base(src)), data, 0o600); err != nil {
		return err
	}

	return nil
}

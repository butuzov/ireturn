package ireturn

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	tests := []testCase{}
	//
	tests = append(tests, testCase{
		name: "Empty Package With No Issues",
		mask: "hello-world*",
	})
	//

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
	want []*analysistest.Result
}

func (tc testCase) test() func(*testing.T) {
	return func(t *testing.T) {
		// -------------------------------------------------------------- setup ----
		dir := t.TempDir()
		t.Cleanup(func() {
			_ = os.RemoveAll(dir)
		})

		if err := tc.xerox(dir); err != nil {
			t.Error(err)
		}

		// --------------------------------------------------------------- test ----
		results := analysistest.Run(&fakeTest{}, dir, NewAnalyzer())

		// ------------------------------------------------------------ results ----
		assert(t, len(tc.want) == len(results[0].Diagnostics), "unexpected results")
	}
}

func (tc testCase) xerox(dest string) error {
	matches, err := filepath.Glob("testdata/" + tc.mask)
	if err != nil {
		return err
	}

	for _, file := range matches {
		if location, err := filepath.Abs(file); err != nil {
			return err
		} else if data, err := ioutil.ReadFile(location); err != nil {
			return err
		} else if err := ioutil.WriteFile(filepath.Join(dest, filepath.Base(file)), data, 0600); err != nil {
			return err
		}
	}

	return nil
}

func assert(t *testing.T, condHappend bool, msg string, args ...interface{}) {
	t.Helper()
	if condHappend {
		return
	}

	t.Errorf(msg, args...)
}

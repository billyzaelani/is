// Package is provides helper function for testing with the standard library.
//
// Comment on the assertion line are used to add a description.
//
// The following failing test:
//
//		func TestFail(t *testing.T) {
// 			// always start tests with this
//			is := is.New(t)
// 			a, b := 1, 2
// 			is.Equal(a, b) // expect to be the same
//		}
//
// Will output:
//
// 		1 != 2 // expect to be the same
package is

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"runtime"
	"strings"
)

// Is is the test helper.
type Is struct {
	t         T
	comments  map[int]string // k:line, v:comment
	arguments map[int]string // k:line, v:argument
}

// New makes a new test helper.
func New(t T) *Is {
	return &Is{t: t}
}

// Equal asserts that a and b are equal.
//
// 		func TestEqual(t *testing.T) {
// 			is := is.New(t)
// 			got := hello("world").
// 			is.Equal(got, "hello world") // greeting the world
// 		}
//
// Will output:
//
// 		wassup world != hello world // greeting the world
func (is *Is) Equal(a, b interface{}) {
	is.t.Helper()
	prefix := "is.Equal"

	if reflect.DeepEqual(a, b) {
		return
	}

	if isNil(a) || isNil(b) {
		is.logf(is.t.Fail, "%s: %s != %s", prefix, valWithType(a), valWithType(b))
		return
	}

	if reflect.ValueOf(a).Type() == reflect.ValueOf(b).Type() {
		is.logf(is.t.Fail, "%s: %v != %v", prefix, a, b)
		return
	}

	is.logf(is.t.Fail, "%s: %s != %s", prefix, valWithType(a), valWithType(b))
}

func (is *Is) logf(failFunc func(), format string, args ...interface{}) {
	is.t.Helper()
	msg := []string{fmt.Sprintf(format, args...)}
	if comment := is.loadComment(); comment != "" {
		msg = append(msg, comment)
	}
	is.t.Log(strings.Join(msg, " "))
	failFunc()
}

func valWithType(v interface{}) string {
	if isNil(v) {
		return "<nil>"
	}
	return fmt.Sprintf("%[1]T(%[1]v)", v)
}

func isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	return false
}

func (is *Is) loadComment() string {
	_, filename, line, _ := runtime.Caller(3) // level of function call to the actual test
	if is.comments == nil {
		is.comments = make(map[int]string)
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, s := range f.Comments {
			line := fset.Position(s.Pos()).Line
			is.comments[line] = "// " + strings.TrimSpace(s.Text())
		}
	}

	return is.comments[line]
}

// NoError assert that err is nil.
//
// 		func TestNoError(t *testing.T) {
//			is := is.New(t)
// 			_, err := findGirlfriend("Anyone?")
// 			is.NoError(err) // poor you
// 		}
//
// Will output:
//
// 		NoError: girlfriend not found // poor you
func (is *Is) NoError(err error) {
	is.t.Helper()
	prefix := "is.NoError"
	if err != nil {
		is.logf(is.t.FailNow, "%s: %q", prefix, err.Error())
	}
}

// Error asserts that err is one of the expectedErrors.
// If no expectedErrors is given, any error will output passed the tests.
func (is *Is) Error(err error, expectedErrors ...error) {
	is.t.Helper()
	prefix := "is.Error"
	if err == nil {
		is.logf(is.t.Fail, "%s: <nil>", prefix)
		return
	}

	lenErr := len(expectedErrors)

	if lenErr == 0 {
		return
	}

	for _, expectedError := range expectedErrors {
		if err == expectedError {
			return
		}
	}

	if lenErr == 1 {
		is.logf(is.t.Fail, "%s: %q != %q", prefix, err.Error(), expectedErrors[0].Error())
		return
	}

	is.logf(is.t.Fail, "%s: %q is not in expected errors", prefix, err.Error())
}

// True asserts that expression is true.
// The expression code itself will be reported if the assertion fails.
//
// 		func TestTrue(t *testing.T) {
// 			is := is.New(t)
// 			val := wallet()
// 			is.True(wallet != 0) // wallet should not be 0
// 		}
//
// Will output:
//
// 		wallet != 0 // wallet should not be 0
func (is *Is) True(expression bool) {
	is.t.Helper()
	prefix := "is.True"

	if expression {
		return
	}

	args := is.loadArgument("True")
	is.logf(is.t.Fail, "%s: %s", prefix, args)
}

func (is *Is) loadArgument(funcName string) string {
	_, filename, line, _ := runtime.Caller(2) // level of function call to the actual test
	if is.arguments == nil {
		is.arguments = make(map[int]string)
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
		if err != nil {
			panic(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			ret, ok := n.(*ast.CallExpr)
			if ok {
				var str strings.Builder
				printer.Fprint(&str, fset, ret)
				if expr := str.String(); strings.Contains(expr, funcName) {
					line := fset.Position(ret.Pos()).Line
					args := strings.ReplaceAll(expr, "\n\t", " ")
					args = args[ret.Lparen-ret.Pos()+1 : len(args)-1]
					is.arguments[line] = args
				}
			}
			return true
		})
	}

	return is.arguments[line]
}

// T is the interface common to testing type.
type T interface {
	Fail()
	FailNow()
	Log(args ...interface{})
	Helper()
}

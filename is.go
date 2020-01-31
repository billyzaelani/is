/*
Package is provides helper function for comparing expected values
to actual values.

Comment add a description

Comment on the assertion lines is optional feature to be printed
as a description in the fail message to keep the API beautifully clean
and easy to use.

The following failing test:

		func TestComment(t *testing.T) {
			is := is.New(t)
			a, b := 1, 2
			is.Equal(a, b) // expect to be the same
		}

Will output:

		is.Equal: 1 != 2 // expect to be the same

Example usage

The example below shows some useful ways to use package is in your test:

		package is_test

		import (
			"errors"
			"strconv"
			"testing"

			"github.com/billyzaelani/is"
		)

		func TestIs(t *testing.T) {
			// always start tests with this
			is := is.New(t)

			i, err := strconv.Atoi("42")

			is.NoError(err)  // passed
			is.Equal(i, 46)  // shouldn't be equal
			is.True(i == 46) // printed the expression code upon failing the test

			j, err = strconv.Atoi("forty two")
			is.Error(err, errors.New("expected errors")) // the error is not expected
			is.NoError(err)                              // we got some error here

			// the code below is not executed because is.NoError uses
			// t.FailNow upon failing the test
			is.True(j)
		}

*/
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

// New makes a new test helper given by T. Any failures will
// reported onto T. Most of the time T will be testing.T from the stdlib.
func New(t T) *Is {
	return &Is{t: t}
}

/*
Equal asserts that a and b are equal. Upon failing the test,
is.Equal also report the data type if a and b has different data type.

		func TestEqual(t *testing.T) {
			is := is.New(t)
			got := hello("girl").
			is.Equal(got, false) // seduce a girl
		}

Will output:

		is.Equal: string(hello girl) != bool(false) // seduce a girl
*/
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

/*
Error asserts that err is one of the expectedErrors.
If no expectedErrors is given, any error will output passed the tests.

		func TestError(t *testing.T) {
			is := is.New(t)
			_, err := findGirlfriend("Anyone?")
			is.Error(err, errors.New("coding")) // its not easy
		}

Will output:

		is.Error: get a girlfriend as programmer? != coding // its not easy
*/
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
		is.logf(is.t.Fail, "%s: %s != %s", prefix, err.Error(), expectedErrors[0].Error())
		return
	}

	is.logf(is.t.Fail, "%s: %s is not in expected errors", prefix, err.Error())
}

/*
NoError assert that err is nil. Upon failing the test,
is.NoError uses t.FailNow so its stop the execution.

		func TestNoError(t *testing.T) {
			is := is.New(t)
			girl, err := findGirlfriend("Anyone?")
			is.NoError(err) // i give up
			is.Equal(girl, nil) // it will not get executed
		}

Will output:

		is.NoError: girlfriend not found // i give up
*/
func (is *Is) NoError(err error) {
	is.t.Helper()
	prefix := "is.NoError"

	if err != nil {
		is.logf(is.t.FailNow, "%s: %s", prefix, err.Error())
	}
}

/*
True asserts that expression is true.
The expression code itself will be reported if the assertion fails.

		func TestTrue(t *testing.T) {
			is := is.New(t)
			money := openTheWallet()
			is.True(money != 0) // money shouldn't be 0 to get a girl
		}

Will output:

		is.True: money != 0 // money shouldn't be 0 to get a girl
*/
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

// PanicFunc is a function to test that function call is panic or not.
type PanicFunc func()

// Panic assert that function f is panic.
func (is *Is) Panic(f PanicFunc) {
	is.t.Helper()

	defer func() {
		is.t.Helper()
		prefix := "is.Panic"

		if recover() != nil {
			return
		}

		is.logf(is.t.Fail, "%s: the function is not panic", prefix)
	}()

	f()
}

// T is the subset of testing.T used by the package is.
type T interface {
	Fail()
	FailNow()
	Log(args ...interface{})
	Helper()
}

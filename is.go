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
func (is *Is) Equal(a, b interface{}) {
	is.t.Helper()
	if reflect.DeepEqual(a, b) {
		return
	}

	if isNil(a) || isNil(b) {
		is.logf("%s != %s", valWithType(a), valWithType(b))
		return
	}

	if reflect.ValueOf(a).Type() == reflect.ValueOf(b).Type() {
		is.logf("%v != %v", a, b)
		return
	}

	is.logf("%s != %s", valWithType(a), valWithType(b))
}

func (is *Is) logf(format string, args ...interface{}) {
	msg := []string{fmt.Sprintf(format, args...)}
	if comment := is.loadComment(); comment != "" {
		msg = append(msg, comment)
	}
	is.t.Error(strings.Join(msg, " "))
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

// NoErr assert that err is nil.
func (is *Is) NoErr(err error) {
	is.t.Helper()
	if err != nil {
		is.logf("err: %s", err.Error())
	}
}

// True asserts that expression is true.
// The expression code itself will be reported if the assertion fails.
func (is *Is) True(expression bool) {
	if expression {
		return
	}

	args := is.loadArgument("True")
	is.logf("false: %s", args)
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
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	Helper()
}

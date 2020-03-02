package is

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func (is *Is) load() {
	is.comments = make(map[string]map[int]string)
	is.arguments = make(map[string]map[int]string)
	_, file, _, _ := runtime.Caller(2)
	root := filepath.Dir(file)
	walkTest := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}

		if strings.HasSuffix(info.Name(), "_test.go") {
			is.comments[path] = loadComment(path)
			is.arguments[path] = loadArgument(path, "True")
		}

		return nil
	}
	filepath.Walk(root, walkTest)
}

func loadComment(path string) map[int]string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	comments := make(map[int]string)
	for _, s := range f.Comments {
		line := fset.Position(s.Pos()).Line
		comments[line] = "// " + strings.TrimSpace(s.Text())
	}
	return comments
}

func loadArgument(path, funcName string) map[int]string {
	arguments := make(map[int]string)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
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
				arguments[line] = args
			}
		}
		return true
	})
	return arguments
}

// logf report the fail depends on failFunc, either t.Fail or t.FailNow.
// skip is how deep the function call to reach the actual test.
func (is *Is) logf(failFunc func(), skip int, format string, args ...interface{}) {
	is.t.Helper()

	msg := []string{fmt.Sprintf(format, args...)}
	if comment := is.loadComment(skip); comment != "" {
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

func (is *Is) loadComment(skip int) string {
	_, file, line, _ := runtime.Caller(skip) // level of function call to the actual test
	return is.comments[file][line]
}

func (is *Is) loadArgument() string {
	_, file, line, _ := runtime.Caller(2) // level of function call to the actual test
	return is.arguments[file][line]
}

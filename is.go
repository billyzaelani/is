package is

import (
	"reflect"
)

// Is is the test helper.
type Is struct {
	t T
}

// New makes a new test helper.
func New(t T) *Is {
	return &Is{t}
}

// Equal asserts that a and b are equal.
func (is *Is) Equal(a, b interface{}) {
	if reflect.DeepEqual(a, b) {
		return
	}

	is.t.Errorf("%v != %v", a, b)
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

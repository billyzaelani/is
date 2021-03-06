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
			"os"
			"strconv"
			"testing"

			"github.com/billyzaelani/is" // neccesary test file will loaded once import the package
		)

		func TestIs(t *testing.T) {
			// always start tests with this
			is := is.New(t)

			i, err := strconv.Atoi("42")

			is.NoError(err)  // passed
			is.Equal(i, 46)  // shouldn't be equal
			is.True(i == 46) // printed the expression code upon failing the test

			j, err = strconv.Atoi("forty two")
			is.Error(err)

			var pathError *os.PathError
			is.ErrorAs(err, &pathError) // err != **os.PathError

			// the code below is not executed because is.ErrorAs uses
			// t.FailNow upon failing the test
			// is.Error and is.NoError also use t.FailNow upon failing the test
			is.True(j)
		}

*/
package is

import (
	"errors"
	"reflect"
)

func init() {
	loadTestFile()
}

var (
	comments  map[string]map[int]string
	arguments map[string]map[int]string
)

// Is is the test helper.
type Is struct {
	T
}

// New makes a new test helper given by T. Any failures will reported onto T.
// Most of the time T will be testing.T from the stdlib.
func New(t T) *Is {
	is := &Is{T: t}
	return is
}

/*
New creates new test helper with the new T.
(In v1.4.1 or later, this function is no different with the New function in package level)

		func TestNew(t *testing.T) {
			is := is.New(t)

			for i := 0; i < 5; i++ {
				t.Run("test"+i, func(t *testing.T) {
					is := is.New(t)
					is.True(true)
				})
			}
		}
*/
func (is *Is) New(t T) *Is {
	return &Is{
		T: t,
	}
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
	if is.T == nil {
		panic("is: T is nil")
	}

	is.Helper()
	prefix := "is.Equal"
	skip := 3

	if reflect.DeepEqual(a, b) {
		return
	}

	if isNil(a) || isNil(b) {
		is.logf(is.T.Fail, skip, "%s: %s != %s", prefix, valWithType(a), valWithType(b))
		return
	}

	if reflect.ValueOf(a).Type() == reflect.ValueOf(b).Type() {
		is.logf(is.Fail, skip, "%s: %v != %v", prefix, a, b)
		return
	}

	is.logf(is.Fail, skip, "%s: %s != %s", prefix, valWithType(a), valWithType(b))
}

/*
Error asserts that err is one of the expectedErrors.
Error uses errors.Is to test the error.
If no expectedErrors is given, any error will output passed the tests.
Error uses t.FailNow upon failing the test.

		func TestError(t *testing.T) {
			is := is.New(t)
			_, err := findGirlfriend("Anyone?")
			is.Error(err, errors.New("coding")) // its not easy
		}

Will output:

		is.Error: get a girlfriend as programmer? != coding // its not easy
*/
func (is *Is) Error(err error, expectedErrors ...error) {
	if is.T == nil {
		panic("is: T is nil")
	}

	is.Helper()
	prefix := "is.Error"
	skip := 3

	if err == nil {
		is.logf(is.FailNow, skip, "%s: <nil>", prefix)
		return
	}

	lenErr := len(expectedErrors)

	if lenErr == 0 {
		return
	}

	for _, expectedError := range expectedErrors {
		if errors.Is(err, expectedError) {
			return
		}
	}

	if lenErr == 1 {
		is.logf(is.FailNow, skip, "%s: %s != %s", prefix, err.Error(), expectedErrors[0].Error())
		return
	}

	is.logf(is.FailNow, skip, "%s: %s != one of the expected errors", prefix, err.Error())
}

/*
ErrorAs asserts that err as target. ErrorAs uses errors.As to test the error.
ErrorAs uses t.FailNow upon failing the test.

		func TestNoError(t *testing.T) {
			is := is.New(t)
			err := errors.New("find a way her heart")
			var pathError *os.PathError
			is.ErrorAs(err, &pathError) // where should I go?
		}

Will output:

		is.ErrorAs: err != **os.PathError // where should I go?
*/
func (is *Is) ErrorAs(err error, target interface{}) {
	if is.T == nil {
		panic("is: T is nil")
	}

	is.Helper()
	prefix := "is.ErrorAs"
	skip := 3

	if !errors.As(err, target) {
		is.logf(is.FailNow, skip, "%s: err != %T", prefix, target)
		return
	}
}

/*
NoError assert that err is nil. NoError uses t.FailNow upon failing the test.

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
	if is.T == nil {
		panic("is: T is nil")
	}

	is.Helper()
	prefix := "is.NoError"
	skip := 3

	if err != nil {
		is.logf(is.FailNow, skip, "%s: %s", prefix, err.Error())
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
	if is.T == nil {
		panic("is: T is nil")
	}

	is.Helper()
	prefix := "is.True"
	skip := 3

	if expression {
		return
	}

	args := is.loadArgument()
	is.logf(is.Fail, skip, "%s: %s", prefix, args)
}

/*
Panic assert that function f is panic.

		func TestPanic(t *testing.T) {
			is := is.New(t)
			panicFunc := func() { panic("single") }
			is.Panic(panicFunc, "really panic", "crazy panic") // ok
		}

Will output:

		is.Panic: single != one of the expected panic values // ok
*/
func (is *Is) Panic(f PanicFunc, expectedValues ...interface{}) {
	if is.T == nil {
		panic("is: T is nil")
	}

	is.Helper()

	defer func(expectedValues ...interface{}) {
		is.Helper()
		prefix := "is.Panic"
		skip := 4

		r := recover()
		if r == nil {
			is.logf(is.Fail, skip, "%s: the function is not panic", prefix)
			return
		}

		lenVal := len(expectedValues)

		if lenVal == 0 {
			return
		}

		for _, v := range expectedValues {
			if reflect.DeepEqual(r, v) {
				return
			}
		}

		if lenVal == 1 {
			is.logf(is.Fail, skip, "%s: %v != %v", prefix, r, expectedValues[0])
			return
		}

		is.logf(is.Fail, skip, "%s: %v != one of the expected panic values", prefix, r)
	}(expectedValues...)

	f()
}

// PanicFunc is a function to test that function call is panic or not.
type PanicFunc func()

// T is the subset of testing.T used by the package is.
type T interface {
	Fail()
	FailNow()
	Log(args ...interface{})
	Helper()
}

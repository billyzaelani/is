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

	_, err = strconv.Atoi("forty two")
	is.Error(err, errors.New("expected errors")) // the error is not expected
	is.NoError(err)                              // we got some error here

	// the code below is not executed because is.NoError is use
	// t.FailNow upon failing the test
	is.True(false)
}

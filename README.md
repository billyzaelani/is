# is [![GoDoc](https://godoc.org/github.com/billyzaelani/is?status.png)](http://godoc.org/github.com/billyzaelani/is) [![Go Report Card](https://goreportcard.com/badge/github.com/billyzaelani/is)](https://goreportcard.com/report/github.com/billyzaelani/is) [![Build Status](https://travis-ci.com/billyzaelani/is.svg?branch=master)](https://travis-ci.com/billyzaelani/is) [![codecov](https://codecov.io/gh/billyzaelani/is/branch/master/graph/badge.svg)](https://codecov.io/gh/billyzaelani/is)

Package is provides helper function for comparing expected values to actual values.

## Comment add a description

Comment on the assertion lines is optional feature to be printed as a description in the fail message to keep the API beautifully clean and easy to use.

The following failing test:

```Go
func TestComment(t *testing.T) {
    is := is.New(t)
    a, b := 1, 2
    is.Equal(a, b) // expect to be the same
}
```

Will output:

```Go
is.Equal: 1 != 2 // expect to be the same
```

## Example usage

The example below shows some useful ways to use package is in your test:

```Go
package is_test

import (
    "errors"
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
```

package is_test

import (
	"errors"
	"fmt"
	"testing"
)

var (
	err1     = errors.New("error 1")
	err2     = errors.New("error 2")
	err3     = errors.New("error 3")
	errWrong = errors.New("something's wrong")
)

type mockT struct {
	state       failState
	msg         string
	helperCount int
}

func (m *mockT) Fail()                   { m.state = fail }
func (m *mockT) FailNow()                { m.state = failNow }
func (m *mockT) Log(args ...interface{}) { m.msg = fmt.Sprint(args...) }
func (m *mockT) Helper()                 { m.helperCount++ }

type failState int

const (
	pass failState = iota
	fail
	failNow
)

func (f failState) String() string {
	var state string
	switch f {
	case pass:
		state = "passed"
	case fail:
		state = "failed"
	case failNow:
		state = "failed now"
	default:
		state = "unknown state"
	}
	return state
}

func assertState(t *testing.T, got, want failState) {
	t.Helper()
	if got != want {
		t.Fatalf("the tests should be %s", want)
	}
}

type QueryError struct{ Query string }

func (e *QueryError) Error() string { return "query: " + e.Query }

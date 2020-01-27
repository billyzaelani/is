package is_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/billyzaelani/is"
)

var errWrong = errors.New("something's wrong")

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
	}
	return state
}

func assertState(t *testing.T, got, want failState) {
	t.Helper()
	if got != want {
		t.Fatalf("the tests should be %s", want)
	}
}

func TestEqual(t *testing.T) {
	prefix := "is.Equal: "
	tests := []struct {
		Name  string
		State failState
		Msg   string
		F     func(is *is.Is)
	}{
		{
			Name:  "equal",
			State: pass,
			Msg:   ``,
			F:     func(is *is.Is) { is.Equal(1, 1) },
		},
		{
			Name:  "not equal",
			State: fail,
			Msg:   prefix + `1 != 2`,
			F:     func(is *is.Is) { is.Equal(1, 2) },
		},
		{
			Name:  "both nil",
			State: pass,
			Msg:   ``,
			F:     func(is *is.Is) { is.Equal(nil, nil) },
		},
		{
			Name:  "different data type",
			State: fail,
			Msg:   prefix + `int(3) != bool(false)`,
			F:     func(is *is.Is) { is.Equal(3, false) },
		},
		{
			Name:  "specific integer",
			State: fail,
			Msg:   prefix + `int32(1) != int64(2)`,
			F:     func(is *is.Is) { is.Equal(int32(1), int64(2)) },
		},
		{
			Name:  "with nil",
			State: fail,
			Msg:   prefix + `<nil> != string(nil)`,
			F:     func(is *is.Is) { is.Equal(nil, "nil") },
		},
		{
			Name:  "nil slice",
			State: fail,
			Msg:   prefix + `[] != [one two]`,
			F: func(is *is.Is) {
				var a []string
				b := []string{"one", "two"}
				is.Equal(a, b)
			},
		},
		{
			Name:  "nil with slice",
			State: fail,
			Msg:   prefix + `<nil> != []string([one two])`,
			F:     func(is *is.Is) { is.Equal(nil, []string{"one", "two"}) },
		},
		{
			Name:  "with comment",
			State: fail,
			Msg:   prefix + `foo != bar // foo is not bar`,
			F: func(is *is.Is) {
				is.Equal("foo", "bar") // foo is not bar
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := new(mockT)
			is := is.New(m)
			tt.F(is)

			assertState(t, m.state, tt.State)
			if m.msg != tt.Msg {
				t.Errorf("%q != %q", m.msg, tt.Msg)
			}
		})
	}
}

func TestNoError(t *testing.T) {
	prefix := "is.NoError: "
	tests := []struct {
		Name  string
		State failState
		Msg   string
		F     func(is *is.Is)
	}{
		{
			Name:  "no error",
			State: pass,
			Msg:   ``,
			F: func(is *is.Is) {
				var err error
				is.NoError(err)
			},
		},
		{
			Name:  "error",
			State: failNow,
			Msg:   prefix + `something's wrong`,
			F:     func(is *is.Is) { is.NoError(errWrong) },
		},
		{
			Name:  "error with comment",
			State: failNow,
			Msg:   prefix + `something's wrong // shouldn't be error`,
			F: func(is *is.Is) {
				is.NoError(errWrong) // shouldn't be error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := new(mockT)
			is := is.New(m)
			tt.F(is)

			assertState(t, m.state, tt.State)
			if m.msg != tt.Msg {
				t.Errorf("got: %s, want: %s", m.msg, tt.Msg)
			}
		})
	}
}

var (
	err1 = errors.New("error 1")
	err2 = errors.New("error 2")
	err3 = errors.New("error 3")
)

func TestError(t *testing.T) {
	prefix := "is.Error: "
	tests := []struct {
		Name  string
		State failState
		Msg   string
		F     func(is *is.Is)
	}{
		{
			Name:  "nil error",
			State: fail,
			Msg:   prefix + `<nil>`,
			F: func(is *is.Is) {
				var err error
				is.Error(err)
			},
		},
		{
			Name:  "nil error with comment",
			State: fail,
			Msg:   prefix + `<nil> // shouldn't be nil`,
			F: func(is *is.Is) {
				var err error
				is.Error(err) // shouldn't be nil
			},
		},
		{
			Name:  "any error",
			State: pass,
			Msg:   ``,
			F:     func(is *is.Is) { is.Error(err1) },
		},
		{
			Name:  "nil with expected error",
			State: fail,
			Msg:   prefix + `<nil>`,
			F: func(is *is.Is) {
				var err error
				is.Error(err, err1)
			},
		},
		{
			Name:  "any error with true expected error",
			State: pass,
			Msg:   ``,
			F: func(is *is.Is) {
				is.Error(err1, err1)
			},
		},
		{
			Name:  "any error with multiple true expected error",
			State: pass,
			Msg:   ``,
			F: func(is *is.Is) {
				is.Error(err2, err1, err2, err3)
			},
		},
		{
			Name:  "any error with false expected error",
			State: fail,
			Msg:   prefix + `error 1 != error 2`,
			F: func(is *is.Is) {
				is.Error(err1, err2)
			},
		},
		{
			Name:  "any error with multiple false expected error",
			State: fail,
			Msg:   prefix + `error 1 is not in expected errors`,
			F: func(is *is.Is) {
				is.Error(err1, err2, err3)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := new(mockT)
			is := is.New(m)
			tt.F(is)

			assertState(t, m.state, tt.State)
			if m.msg != tt.Msg {
				t.Errorf("got: %s, want: %s", m.msg, tt.Msg)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	prefix := "is.True: "
	tests := []struct {
		Name  string
		State failState
		Msg   string
		F     func(is *is.Is)
	}{
		{
			Name:  "true",
			State: pass,
			Msg:   ``,
			F: func(is *is.Is) {
				is.True(1 == 1) // true
			},
		},
		{
			Name:  "false",
			State: fail,
			Msg:   prefix + `1 == 2 // comment`,
			F: func(is *is.Is) {
				is.True(1 == 2) // comment
			},
		},
		{
			Name:  "extra parentheses",
			State: fail,
			Msg:   prefix + `(1 == 2) // comment`,
			F: func(is *is.Is) {
				is.True((1 == 2)) // comment
			},
		},
		{
			Name:  "new line",
			State: fail,
			Msg:   prefix + `(1 == 2) && false`,
			F: func(is *is.Is) {
				is.True((1 == 2) &&
					false)
			},
		},
		{
			Name:  "multi line",
			State: fail,
			Msg:   prefix + `(1 == 2) && false || false`,
			F: func(is *is.Is) {
				is.True((1 == 2) &&
					false ||
					false)
			},
		},
		{
			Name:  "multi line with comment in first line",
			State: fail,
			Msg:   prefix + `(1 == 2) && false || false // comment`,
			F: func(is *is.Is) {
				is.True((1 == 2) && // comment
					false ||
					false)
			},
		},
		{
			Name:  "multi line with comment in non-first line",
			State: fail,
			Msg:   prefix + `(1 == 2) && false || false`,
			F: func(is *is.Is) {
				is.True((1 == 2) &&
					false || // cannot be printed
					false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := new(mockT)
			is := is.New(m)
			tt.F(is)

			assertState(t, m.state, tt.State)
			if m.msg != tt.Msg {
				t.Errorf("%q != %q", m.msg, tt.Msg)
			}
		})
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		Name string
		F    func(is *is.Is)
		Want int
	}{
		{
			Name: "Equal",
			F:    func(is *is.Is) { is.Equal(1, 2) },
			Want: 2,
		},
		{
			Name: "NoError",
			F:    func(is *is.Is) { is.NoError(errWrong) },
			Want: 2,
		},
		{
			Name: "Error",
			F:    func(is *is.Is) { is.Error(nil) },
			Want: 2,
		},
		{
			Name: "True",
			F:    func(is *is.Is) { is.True(1 == 2) },
			Want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := new(mockT)
			is := is.New(m)
			tt.F(is)

			if m.helperCount != tt.Want {
				t.Errorf("%d != %d", m.helperCount, tt.Want)
			}
		})
	}
}

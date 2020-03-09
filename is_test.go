package is_test

import (
	"errors"
	"testing"

	assert "github.com/billyzaelani/is"
)

// load the test file upfront with nil T
var is = assert.New(nil)

func TestEqual(t *testing.T) {
	prefix := "is.Equal: "
	tests := []struct {
		desc  string
		state failState
		msg   string
		f     func(is *assert.Is)
	}{
		{"equal", pass, ``,
			func(is *assert.Is) { is.Equal(1, 1) }},
		{"not equal", fail, prefix + `1 != 2`,
			func(is *assert.Is) { is.Equal(1, 2) }},
		{"both nil", pass, ``,
			func(is *assert.Is) { is.Equal(nil, nil) }},
		{"different data type", fail, prefix + `int(3) != bool(false)`,
			func(is *assert.Is) { is.Equal(3, false) }},
		{"specific integer", fail, prefix + `int32(1) != int64(2)`,
			func(is *assert.Is) { is.Equal(int32(1), int64(2)) }},
		{"with nil", fail, prefix + `<nil> != string(nil)`,
			func(is *assert.Is) { is.Equal(nil, "nil") }},
		{"nil slice", fail, prefix + `[] != [one two]`,
			func(is *assert.Is) { is.Equal([]string{}, []string{"one", "two"}) }},
		{"nil with slice", fail, prefix + `<nil> != []string([one two])`,
			func(is *assert.Is) { is.Equal(nil, []string{"one", "two"}) }},
		{"with comment", fail, prefix + `foo != bar // foo is not bar`,
			func(is *assert.Is) { is.Equal("foo", "bar") /* foo is not bar */ }},
	}

	for _, tt := range tests {
		// fix for-range local variable when using t.Parallel that reduces code coverage
		// see: https://gist.github.com/posener/92a55c4cd441fc5e5e85f27bca008721
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			assertState(t, m.state, tt.state)
			if m.msg != tt.msg {
				t.Errorf("%q != %q", m.msg, tt.msg)
			}
		})
	}
}

func TestNoError(t *testing.T) {
	prefix := "is.NoError: "
	tests := []struct {
		name  string
		state failState
		msg   string
		f     func(is *assert.Is)
	}{
		{"no error", pass, ``,
			func(is *assert.Is) { is.NoError(nil) }},
		{"error", failNow, prefix + `something's wrong`,
			func(is *assert.Is) { is.NoError(errWrong) }},
		{"error with comment", failNow, prefix + `something's wrong // shouldn't be error`,
			func(is *assert.Is) { is.NoError(errWrong) /* shouldn't be error*/ }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			assertState(t, m.state, tt.state)
			if m.msg != tt.msg {
				t.Errorf("got: %s, want: %s", m.msg, tt.msg)
			}
		})
	}
}

func TestError(t *testing.T) {
	prefix := "is.Error: "
	tests := []struct {
		name  string
		state failState
		msg   string
		f     func(is *assert.Is)
	}{
		{"nil error", failNow, prefix + `<nil>`,
			func(is *assert.Is) { is.Error(nil) }},
		{"nil error with comment", failNow, prefix + `<nil> // shouldn't be nil`,
			func(is *assert.Is) { is.Error(nil) /* shouldn't be nil*/ }},
		{"any error", pass, ``,
			func(is *assert.Is) { is.Error(err1) }},
		{"nil with expected error", failNow, prefix + `<nil>`,
			func(is *assert.Is) { is.Error(nil, err1) }},
		{"any error with true expected error", pass, ``,
			func(is *assert.Is) { is.Error(err1, err1) }},
		{"any error with multiple true expected error", pass, ``,
			func(is *assert.Is) { is.Error(err2, err1, err2, err3) }},
		{"any error with false expected error", failNow, prefix + `error 1 != error 2`,
			func(is *assert.Is) { is.Error(err1, err2) }},
		{"any error with multiple false expected error", failNow, prefix + `error 1 != one of the expected errors`,
			func(is *assert.Is) { is.Error(err1, err2, err3) }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			assertState(t, m.state, tt.state)
			if m.msg != tt.msg {
				t.Errorf("got: %s, want: %s", m.msg, tt.msg)
			}
		})
	}
}

func TestErrorAs(t *testing.T) {
	prefix := "is.ErrorAs: "
	tests := []struct {
		name  string
		state failState
		msg   string
		f     func(is *assert.Is)
	}{
		{"pass", pass, ``,
			func(is *assert.Is) {
				var e *QueryError
				is.ErrorAs(&QueryError{"SELECT column_name(s) FROM table_name"}, &e)
			}},
		{"fail", failNow, prefix + `err != **is_test.QueryError // it's something else`,
			func(is *assert.Is) {
				var e *QueryError
				is.ErrorAs(errors.New("it's not query error"), &e) // it's something else
			}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			assertState(t, m.state, tt.state)
			if m.msg != tt.msg {
				t.Errorf("%q != %q", m.msg, tt.msg)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	prefix := "is.True: "
	tests := []struct {
		name  string
		state failState
		msg   string
		f     func(is *assert.Is)
	}{
		{"true", pass, ``,
			func(is *assert.Is) { is.True(1 == 1) }},
		{"false", fail, prefix + `1 == 2 // false`,
			func(is *assert.Is) { is.True(1 == 2) /* false*/ }},
		{"extra parentheses", fail, prefix + `(1 == 2) // comment`,
			func(is *assert.Is) { is.True((1 == 2)) /* comment */ }},
		{"new line", fail, prefix + `(1 == 2) && false`,
			func(is *assert.Is) {
				is.True((1 == 2) &&
					false)
			}},
		{"multi line", fail, prefix + `(1 == 2) && false || false`,
			func(is *assert.Is) {
				is.True((1 == 2) &&
					false ||
					false)
			}},
		{"multi line with comment in first line", fail, prefix + `(1 == 2) && false || false // comment`,
			func(is *assert.Is) {
				is.True((1 == 2) && // comment
					false ||
					false)
			}},
		{"multi line with comment in non-first line", fail, prefix + `(1 == 2) && false || false`,
			func(is *assert.Is) {
				is.True((1 == 2) &&
					false || // cannot be printed
					false)
			}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			assertState(t, m.state, tt.state)
			if m.msg != tt.msg {
				t.Errorf("%q != %q", m.msg, tt.msg)
			}
		})
	}
}

func TestPanic(t *testing.T) {
	prefix := "is.Panic: "
	tests := []struct {
		name  string
		state failState
		msg   string
		f     func(is *assert.Is)
	}{
		{"panic", pass, ``,
			func(is *assert.Is) {
				panicFunc := func() { panic("i'm panic") }
				is.Panic(panicFunc)
			}},
		{"not panic", fail, prefix + `the function is not panic`,
			func(is *assert.Is) {
				calmFunc := func() { _ = "i'm calm" }
				is.Panic(calmFunc)
			}},
		{"not panic with comment", fail, prefix + `the function is not panic // with comment`,
			func(is *assert.Is) {
				calmFunc := func() { _ = "i'm calm" }
				is.Panic(calmFunc) // with comment
			}},
		{"panic with true panic value", pass, ``,
			func(is *assert.Is) {
				panicFunc := func() { panic("i'm panic") }
				is.Panic(panicFunc, "i'm panic", "is this panic")
			}},
		{"panic with false panic value", fail, prefix + `i'm panic != are you panic`,
			func(is *assert.Is) {
				panicFunc := func() { panic("i'm panic") }
				is.Panic(panicFunc, "are you panic")
			}},
		{"panic with multiple false panic value", fail, prefix + `i'm panic != one of the expected panic values`,
			func(is *assert.Is) {
				panicFunc := func() { panic("i'm panic") }
				is.Panic(panicFunc, "are you panic", "are you crazy")
			}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			assertState(t, m.state, tt.state)
			if m.msg != tt.msg {
				t.Errorf("%q != %q", m.msg, tt.msg)
			}
		})
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		name string
		want int
		f    func(is *assert.Is)
	}{
		{"Equal", 2, func(is *assert.Is) { is.Equal(1, 2) }},
		{"NoError", 2, func(is *assert.Is) { is.NoError(errWrong) }},
		{"Error", 2, func(is *assert.Is) { is.Error(nil) }},
		{"ErrorAs", 2, func(is *assert.Is) {
			var e *QueryError
			is.ErrorAs(errors.New("it's not query error"), &e)
		}},
		{"True", 2, func(is *assert.Is) { is.True(1 == 2) }},
		{"Panic", 3, func(is *assert.Is) { is.Panic(func() {}) }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mockT)
			is := is.New(m)
			tt.f(is)

			if m.helperCount != tt.want {
				t.Errorf("%d != %d", m.helperCount, tt.want)
			}
		})
	}
}

func TestHelperPanic(t *testing.T) {
	tests := []struct {
		name string
		f    func()
	}{
		{"is.Equal panic", func() { is.Equal(1, 1) }},
		{"is.NoError panic", func() { is.NoError(nil) }},
		{"is.Error panic", func() { is.Error(nil) }},
		{"is.ErrorAs panic", func() { is.ErrorAs(nil, nil) }},
		{"is.True panic", func() { is.True(false) }},
		{"is.Panic panic", func() { is.Panic(nil) }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				errMsg := "is: T is nil"
				if err := recover(); err != errMsg {
					t.Errorf("%q != %q", err, errMsg)
				}
			}()

			tt.f()
		})
	}
}

func TestExportedField(t *testing.T) {
	t.Run("Fail", func(t *testing.T) {
		m := new(mockT)
		assert.New(m).Fail()
		if m.state != fail {
			t.Errorf("%q != %q", m.state, fail)
		}
	})
	t.Run("FailNow", func(t *testing.T) {
		m := new(mockT)
		assert.New(m).FailNow()
		if m.state != failNow {
			t.Errorf("%q != %q", m.state, fail)
		}
	})
	t.Run("Helper", func(t *testing.T) {
		m := new(mockT)
		assert.New(m).Helper()
		if m.helperCount != 1 {
			t.Errorf("%q != %d", m.helperCount, 1)
		}
	})
}

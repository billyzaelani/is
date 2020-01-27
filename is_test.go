package is_test

import (
	"errors"
	"fmt"
	"github.com/billyzaelani/is"
	"testing"
)

type mockT struct {
	fail        bool
	failNow     bool
	message     string
	helperCount int
}

func (m *mockT) Fail()                   { m.fail = true }
func (m *mockT) FailNow()                { m.failNow = true }
func (m *mockT) Log(args ...interface{}) { m.message = fmt.Sprint(args...) }
func (m *mockT) Helper()                 { m.helperCount++ }
func (m *mockT) out() string             { return m.message }

func TestEqual(t *testing.T) {
	tests := []struct {
		Name string
		Got  func(is *is.Is)
		Want string
	}{
		{
			Name: "equal",
			Got:  func(is *is.Is) { is.Equal(1, 1) },
			Want: ``,
		},
		{
			Name: "not equal",
			Got:  func(is *is.Is) { is.Equal(1, 2) },
			Want: `1 != 2`,
		},
		{
			Name: "both nil",
			Got:  func(is *is.Is) { is.Equal(nil, nil) },
			Want: ``,
		},
		{
			Name: "different data type",
			Got:  func(is *is.Is) { is.Equal(3, false) },
			Want: `int(3) != bool(false)`,
		},
		{
			Name: "specific integer",
			Got:  func(is *is.Is) { is.Equal(int32(1), int64(2)) },
			Want: `int32(1) != int64(2)`,
		},
		{
			Name: "with nil",
			Got:  func(is *is.Is) { is.Equal(nil, "nil") },
			Want: `<nil> != string(nil)`,
		},
		{
			Name: "nil slice",
			Got: func(is *is.Is) {
				var a []string
				b := []string{"one", "two"}
				is.Equal(a, b)
			},
			Want: `[] != [one two]`,
		},
		{
			Name: "nil with slice",
			Got:  func(is *is.Is) { is.Equal(nil, []string{"one", "two"}) },
			Want: `<nil> != []string([one two])`,
		},
		{
			Name: "with comment",
			Got: func(is *is.Is) {
				is.Equal("foo", "bar") // foo is not bar
			},
			Want: "foo != bar // foo is not bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := &mockT{}
			is := is.New(m)
			tt.Got(is)
			got := m.out()
			want := tt.Want

			if got != want {
				t.Errorf("%q != %q", got, want)
			}
		})
	}
}

func TestNoErr(t *testing.T) {
	tests := []struct {
		Name string
		Got  func(is *is.Is)
		Want string
	}{
		{
			Name: "no error",
			Got: func(is *is.Is) {
				var err error
				is.NoErr(err)
			},
			Want: ``,
		},
		{
			Name: "error",
			Got:  func(is *is.Is) { is.NoErr(errors.New("something's wrong")) },
			Want: `err: something's wrong`,
		},
		{
			Name: "error with comment",
			Got: func(is *is.Is) {
				is.NoErr(errors.New("something's wrong")) // shouldn't be error
			},
			Want: `err: something's wrong // shouldn't be error`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := &mockT{}
			is := is.New(m)
			tt.Got(is)
			got := m.out()
			want := tt.Want

			if got != want {
				t.Errorf("%q != %q", got, want)
			}
		})
	}
}

func TestTrue(t *testing.T) {
	tests := []struct {
		Name string
		Got  func(is *is.Is)
		Want string
	}{
		{
			Name: "true",
			Got: func(is *is.Is) {
				is.True(1 == 1) // true
			},
			Want: ``,
		},
		{
			Name: "false",
			Got: func(is *is.Is) {
				is.True(1 == 2) // comment
			},
			Want: `false: 1 == 2 // comment`,
		},
		{
			Name: "extra parentheses",
			Got: func(is *is.Is) {
				is.True((1 == 2)) // comment
			},
			Want: `false: (1 == 2) // comment`,
		},
		{
			Name: "new line",
			Got: func(is *is.Is) {
				is.True((1 == 2) &&
					false)
			},
			Want: `false: (1 == 2) && false`,
		},
		{
			Name: "multi line",
			Got: func(is *is.Is) {
				is.True((1 == 2) &&
					false ||
					false)
			},
			Want: `false: (1 == 2) && false || false`,
		},
		{
			Name: "multi line with comment in first line",
			Got: func(is *is.Is) {
				is.True((1 == 2) && // comment
					false ||
					false)
			},
			Want: `false: (1 == 2) && false || false // comment`,
		},
		{
			Name: "multi line with comment in non-first line",
			Got: func(is *is.Is) {
				is.True((1 == 2) &&
					false || // cannot be printed
					false)
			},
			Want: `false: (1 == 2) && false || false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := &mockT{}
			is := is.New(m)
			tt.Got(is)
			got := m.out()
			want := tt.Want

			if got != want {
				t.Errorf("%q != %q", got, want)
			}
		})
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		Name string
		Got  func(is *is.Is)
		Want int
	}{
		{
			Name: "Equal",
			Got:  func(is *is.Is) { is.Equal(1, 2) },
			Want: 2,
		},
		{
			Name: "NoErr",
			Got:  func(is *is.Is) { is.NoErr(errors.New("something's wrong")) },
			Want: 2,
		},
		{
			Name: "True",
			Got:  func(is *is.Is) { is.True(1 == 2) },
			Want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := &mockT{}
			is := is.New(m)
			tt.Got(is)
			got := m.helperCount
			want := tt.Want

			if got != want {
				t.Errorf("%d != %d", got, want)
			}
		})
	}
}

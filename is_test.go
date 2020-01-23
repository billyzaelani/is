package is

import (
	"fmt"
	"testing"
)

type mockT struct {
	fail bool
	skip bool
	s    string
}

func (m *mockT) Error(args ...interface{})                 { m.Log(args...); m.Fail() }
func (m *mockT) Errorf(format string, args ...interface{}) { m.Logf(format, args...); m.Fail() }
func (m *mockT) Fail()                                     { m.fail = true }
func (m *mockT) FailNow()                                  { m.fail = true }
func (m *mockT) Failed() bool                              { return m.fail }
func (m *mockT) Fatal(args ...interface{})                 { m.Log(args...); m.FailNow() }
func (m *mockT) Fatalf(format string, args ...interface{}) { m.Logf(format, args...); m.FailNow() }
func (m *mockT) Log(args ...interface{})                   { m.s = fmt.Sprint(args...) }
func (m *mockT) Logf(format string, args ...interface{})   { m.s = fmt.Sprintf(format, args...) }
func (m *mockT) Name() string                              { return "" }
func (m *mockT) Skip(args ...interface{})                  { m.Log(args...); m.SkipNow() }
func (m *mockT) SkipNow()                                  { m.skip = true }
func (m *mockT) Skipf(format string, args ...interface{})  { m.Logf(format, args...); m.SkipNow() }
func (m *mockT) Skipped() bool                             { return m.skip }
func (m *mockT) Helper()                                   {}
func (m *mockT) out() string                               { return m.s }

func TestEqual(t *testing.T) {
	tests := []struct {
		Name string
		Got  func(is *Is)
		Want string
	}{
		{
			Name: "equal",
			Got:  func(is *Is) { is.Equal(1, 1) },
			Want: ``,
		},
		{
			Name: "not equal",
			Got:  func(is *Is) { is.Equal(1, 2) },
			Want: `1 != 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := &mockT{}
			is := New(m)
			tt.Got(is)
			got := m.out()
			want := tt.Want

			if got != want {
				t.Errorf("%q != %q", got, want)
			}
		})
	}
}

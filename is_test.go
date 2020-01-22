package is

import (
	"fmt"
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

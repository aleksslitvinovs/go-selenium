package types

// TestFunction describes one test for the given session. It is used in
// selenium.Run() to execute tests.
type TestFunction func(s Sessioner)

package tui

// Messages for internal events

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type publishSuccessMsg struct{}
type publishErrorMsg struct{ err error }

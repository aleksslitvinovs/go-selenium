package types

type Clienter interface {
	CreateSession() (Sessioner, error)
	MustStop()
	Stop() error
	RaiseErrors()
}

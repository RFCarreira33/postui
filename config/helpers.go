package config

type Status int

const (
	NUM_TABS        = 5
	Base     Status = iota
	Params
	Auth
	Headers
	Body
)

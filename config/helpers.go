package config

type Status int

const NUM_TABS = 5

const (
	// Tabs
	Base Status = iota
	Params
	Auth
	Headers
	Body
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

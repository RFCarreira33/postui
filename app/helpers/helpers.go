package helpers

const NUM_TABS = 5

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// General models

type Tab int

const (
	// Tabs
	Base Tab = iota
	Params
	Auth
	Headers
	Body
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Visual
)

func (m Mode) IsNormal() bool {
	return m == Normal
}

func (m Mode) IsInsert() bool {
	return m == Insert
}

func (m Mode) IsVisual() bool {
	return m == Visual
}

package flags

import "flag"

type TestMode int

const (
	FuncTest TestMode = iota
	PerfTest
)

var (
	ModeFlag = flag.Int("testmode", int(FuncTest), "Test mode: 0-functional, 1-performance")
)

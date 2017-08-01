package internal

import (
	"os"
	"runtime/pprof"

	"gomatcha.io/bridge"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/internal printStack", printStack)
}

func printStack() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

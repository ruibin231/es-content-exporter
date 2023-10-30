package async

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				const call = 2
				_, file, line, ok := runtime.Caller(call)
				if !ok {
					file = "???"
					line = 0
				}
				fmt.Printf("\npanic err: %s:%d %s\n\n", file, line, err)
				debug.PrintStack()
			}
		}()
		f()
	}()
}

package sws

import (
	"fmt"
	"os"
	"time"
)

type Logger func(...interface{})

var (
	DebugLog Logger = func(v ...interface{}) {}
	ErrorLog Logger = func(v ...interface{}) {
		fmt.Fprintf(os.Stderr, "[%s] ", time.Now().Format(time.RFC3339))
		fmt.Fprintln(os.Stderr, v...)
	}
)

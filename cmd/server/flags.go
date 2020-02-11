package main

import (
	"flag"
	"os"
)

func stringFlag(long, short, def, envvar, desc string) *string {
	out := flag.String(long, def, desc)
	if short != "" {
		flag.StringVar(out, short, def, desc)
	}
	if envvar != "" {
		if v := os.Getenv(envvar); v != "" {
			def = v
		}
	}
	return out
}

func boolFlag(long, short string, def bool, envvar, desc string) *bool {
	if envvar != "" {
		if v := os.Getenv(envvar); v != "" {
			def = true
		}
	}
	out := flag.Bool(long, def, desc)
	if short != "" {
		flag.BoolVar(out, short, def, desc)
	}
	return out
}

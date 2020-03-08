package main

import (
	"flag"
	"os"
)

func stringFlag(long, short, def, envvar, desc string) *string {
	if envvar != "" {
		if v := os.Getenv(envvar); v != "" {
			def = v
		}
	}
	out := flag.String(long, def, desc)
	if short != "" {
		flag.StringVar(out, short, def, desc)
	}
	if out != nil && *out == "" {
		return nil
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

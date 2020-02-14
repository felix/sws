package main

import (
	"fmt"
	"html/template"
	"time"
)

var funcMap = template.FuncMap{
	"sparkline": func(id int) string {
		// This will enable "caching" for an hour
		now := time.Now().Truncate(time.Hour)
		//then := now.Add(-720 * time.Hour)
		then := now.Add(-24 * time.Hour)
		return fmt.Sprintf("/sites/%d/sparklines/%d-%d.svg", id, then.Unix(), now.Unix())
	},
}

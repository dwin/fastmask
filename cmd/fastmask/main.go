package main

import "github.com/dwin/fastmask/internal/cli"

var (
	date    string
	commit  string
	version string
)

func main() {
	// nolint:errcheck
	cli.LoadFastmask(date, version, commit).Execute()
}

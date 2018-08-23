package main

import (
	"os"

	"github.com/covexo/devspace/cmd"
	"github.com/covexo/devspace/pkg/devspace/upgrade"
)

var version string

func main() {
	upgrade.SetVersion(version)

	cmd.Execute()
	os.Exit(0)
}

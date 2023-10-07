package main

import (
	_ "embed"

	"github.com/epsxy/flower/pkg/cmd"
)

//go:embed VERSION
var version string

func main() {
	cmd.Execute(version)
}

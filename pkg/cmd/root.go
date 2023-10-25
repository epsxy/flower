package cmd

import (
	"fmt"
	"os"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/log"
	"github.com/spf13/cobra"
)

// default level: error
var logLevel = log.LogLevelError

var root = &cobra.Command{
	Use:     "flower",
	Short:   "Flower SQL-dump to PlantUML ERD",
	Long:    "A Go program to help you parse SQL dumps in PlantUML ERD",
	Version: "v0.0.0", // will be filled in `Execute` entrypoint
}

func SetGlobalFlags(cmd *cobra.Command) {
	// v, err := cmd.Flags().GetBool("verbose")
	// if err != nil {
	// 	v = false
	// }
	// global.SetVerbose(v)

	// d, err := cmd.Flags().GetBool("dry-run")
	// if err != nil {
	// 	d = false
	// }
	// global.SetDryRun(d)

	// set log level from flag
	global.SetLogger(logLevel.ToSlogLevel())
}

// Execute is the root entrypoint of the Cobra CLI
func Execute(version string) {
	// setup version
	root.Version = version
	// add commands and flags
	Parse.PersistentFlags().String("input", "", "Path to SQL file to read")
	Parse.PersistentFlags().String("output", "", "Path to PlantUML file to write (including '.plantuml' extension)")
	Parse.PersistentFlags().Bool("partition", false, "Split the disconnected data graph into connected data graph, one per file. Default=false")
	root.AddCommand(Parse)

	root.PersistentFlags().Var(&logLevel, "log-level", "Log level: 'debug', 'info', 'debug' or 'error'")
	//root.PersistentFlags().Bool("dry-run", false, "enable dry-run mode")
	//root.PersistentFlags().Bool("verbose", false, "enable verbose mode")

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

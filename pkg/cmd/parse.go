package cmd

import (
	"fmt"
	"os"

	"github.com/epsxy/flower/pkg/reader"
	"github.com/spf13/cobra"
)

// Parse
var Parse = &cobra.Command{
	Use:   "parse",
	Short: "Run parse command ",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		SetGlobalFlags(cmd)
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")

		dat, err := os.ReadFile(input)
		if err != nil {
			os.Exit(1)
		}
		tree := reader.Read(string(dat))
		err = os.WriteFile(output, []byte(tree.Build()), 0644)
		if err != nil {
			fmt.Println("unable to write file")
			os.Exit(1)
		}
	},
}

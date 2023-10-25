package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/epsxy/flower/pkg/global"
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
		partition, _ := cmd.Flags().GetBool("partition")

		logger := global.GetLogger()

		dat, err := os.ReadFile(input)
		if err != nil {
			os.Exit(1)
		}
		tree := reader.Read(string(dat))

		if !partition {
			res := tree.Build()
			err := os.WriteFile(output, []byte(res), 0644)
			if err != nil {
				logger.Error("unable to save file", "filename", output)
			}
		} else {
			res := tree.BuildWithPartitions()
			for i, r := range res {
				var prefix, extension string
				split := strings.SplitN(output, ".", 2)
				if len(split) == 2 {
					prefix = split[0]
					extension = split[1]
				} else {
					prefix = output
					extension = ".plantuml"
				}
				filename := fmt.Sprintf("%s_%d.%s", prefix, i, extension)
				err := os.WriteFile(filename, []byte(r), 0644)
				if err != nil {
					logger.Error("unable to save file", "filename", filename)
				}
			}
		}
		if err != nil {
			fmt.Println("unable to write file")
			os.Exit(1)
		}
	},
}

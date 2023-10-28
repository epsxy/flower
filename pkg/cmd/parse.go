package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/epsxy/flower/pkg/global"
	"github.com/epsxy/flower/pkg/model"
	"github.com/epsxy/flower/pkg/reader"
	"github.com/epsxy/flower/pkg/writer"
	"github.com/spf13/cobra"
)

// Parse
var Parse = &cobra.Command{
	Use:   "parse",
	Short: "Run parse command ",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		/*
			// mark distance required when any of these flag is provided
			maxPartitionSize, _ := cmd.Flags().GetInt("max-partition")
			weightEdge, _ := cmd.Flags().GetInt("weight-edge")
			weightDistance, _ := cmd.Flags().GetInt("weight-distance")
			if maxPartitionSize != 0 {
				cmd.MarkFlagRequired("distance")
			}
			if weightEdge != 0 {
				cmd.MarkFlagRequired("distance")
			}
			if weightDistance != 0 {
				cmd.MarkFlagRequired("distance")
			}
			// guarantee that distance is not empty when provided and matches expected values
			dist, _ := cmd.Flags().GetString("distance")
			var distance model.DistanceNorm
			if dist != "" {
				err := distance.Set(dist)
				if err != nil {
					fmt.Println("invalid distance parameter")
					os.Exit(1)
				}
			}
		*/
	},
	Run: func(cmd *cobra.Command, args []string) {
		SetGlobalFlags(cmd)
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		splitUnconnected, _ := cmd.Flags().GetBool("split-unconnected")
		dist, _ := cmd.Flags().GetString("distance")
		maxPartitionSize, _ := cmd.Flags().GetInt("max-partition")
		weightEdge, _ := cmd.Flags().GetInt("weight-edge")
		weightDistance, _ := cmd.Flags().GetInt("weight-distance")

		// set distance
		var distance model.DistanceNorm
		distance.Set(dist)

		// assign build options from input
		options := &model.UMLTreeOptions{
			SplitUnconnected: splitUnconnected,
			SplitDistance:    distance != "",
			DistanceNorm:     distance,
			MaxPartitionSize: maxPartitionSize,
			WeightEdge:       weightEdge,
			WeightDistance:   weightDistance,
		}

		logger := global.GetLogger()

		logger.Warn("distance", "distance", distance)
		logger.Warn("args", "args", args)

		dat, err := os.ReadFile(input)
		if err != nil {
			os.Exit(1)
		}

		tree := reader.Read(string(dat))

		res := writer.Build(tree.SetOptions(options))
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
	},
}

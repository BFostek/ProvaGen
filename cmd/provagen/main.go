package main

import (
	"fmt"
	"os"

	"github.com/BFostek/ProvaGen/pkg/generator"
	"github.com/spf13/cobra"
)

var (
	outputDir string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "provagen [problem-name]",
		Short: "ProvaGen is a tool for generating challenge files",
		Long:  `ProvaGen creates the necessary files and structure for coding challenges`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			problemName := args[0]
			
			gen := generator.GeneratorFactory()
			if gen == nil {
				fmt.Println("Invalid generator")
				os.Exit(1)
			}
			
			if err := gen.Generate(outputDir, problemName); err != nil {
				fmt.Printf("Error generating files: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Successfully generated files for '%s' in %s\n", problemName, outputDir)
		},
	}

	// Set up flags
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "/home/breno/dev/challenges", "output directory for generated files")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

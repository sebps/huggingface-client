package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huggingface-cli",
	Short: "CLI to manage Hugging Face endpoints",
	Long:  `A Cobra CLI to manage Hugging Face hosted inference endpoints.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

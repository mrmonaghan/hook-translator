package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "stitch",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("template-dir", "./templates", "specify path to template directory")
	rootCmd.PersistentFlags().String("rules-dir", "./rules", "specify path to rules directory")
	rootCmd.PersistentFlags().Bool("exit-on-parse-error", false, "")
}

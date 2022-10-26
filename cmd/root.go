package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var (
	version    = "0.1"
	serverAddr string
)

var rootCmd = &cobra.Command{
	Use:   path.Base(os.Args[0]),
	Short: "Makabaka Api Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: ", version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

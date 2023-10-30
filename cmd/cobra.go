package cmd

import (
	"es-content-export/cmd/service"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:               "exporter",
		Short:             "exporter",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Args:              nil,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("欢迎使用, 可以使用 -h 查看命令帮助!")
		},
	}
)

func init() {
	rootCmd.AddCommand(service.StartServerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

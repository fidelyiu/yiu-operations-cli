package cmd

import "yiu-ops/internal/version"

func init() {
	rootCmd.Version = version.BuildVersionOutput()
	rootCmd.AddCommand(version.NewCommand())
}

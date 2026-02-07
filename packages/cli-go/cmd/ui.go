package cmd

import "yiu-ops/internal/ui"

func init() {
	rootCmd.AddCommand(ui.NewCommand(appCtx))
}

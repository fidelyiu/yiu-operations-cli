package cmd

import "yiu-ops/internal/docs"

func init() {
	rootCmd.AddCommand(docs.NewCommand(appCtx))
}

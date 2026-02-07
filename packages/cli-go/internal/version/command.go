package version

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "打印当前 CLI 版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(BuildVersionOutput())
		},
	}

	return cmd
}

func BuildVersionOutput() string {
	versionValue, commitValue, buildDateValue := readVersionInfo()
	return fmt.Sprintf("version: %s\ncommit: %s\nbuildDate: %s", versionValue, commitValue, buildDateValue)
}

func readVersionInfo() (string, string, string) {
	versionValue := Version
	commitValue := Commit
	buildDateValue := BuildDate

	if info, ok := debug.ReadBuildInfo(); ok {
		if versionValue == "" || versionValue == "dev" {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				versionValue = info.Main.Version
			}
		}

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if commitValue == "none" {
					commitValue = setting.Value
				}
			case "vcs.time":
				if buildDateValue == "unknown" {
					buildDateValue = setting.Value
				}
			}
		}
	}

	return versionValue, commitValue, buildDateValue
}

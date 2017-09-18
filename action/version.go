package action

import (
	"fmt"                                             // go-core
	"github.com/roscopecoltran/sniperkit-limo/config" // app-config
	"github.com/sirupsen/logrus"                      // logs-logrus
	"github.com/spf13/cobra"                          // cli-cmd
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

// VersionCmd shows the version
var VersionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Display version information",
	Long:    fmt.Sprintf("Display version information for %s.", config.ProgramName),
	Example: fmt.Sprintf("  %s version", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		getOutput().Info(config.Version)
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"prefix":      "app-action",
			"src.file":    "action/version.go",
			"cmd.name":    "VersionCmd",
			"method.name": "init()",
			"var.options": options,
		}).Info("registering command...")
	RootCmd.AddCommand(VersionCmd)
}

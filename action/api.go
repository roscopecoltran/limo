package action

import (
	"fmt"																					// go-core
	"github.com/roscopecoltran/sniperkit-limo/config" 										// app-config
	//"github.com/roscopecoltran/sniperkit-limo/service" 									// svc-registry
	//"github.com/roscopecoltran/sniperkit-limo/model" 										// data-models
	"github.com/spf13/cobra" 																// cli-cmd
	"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

// ClassifyCmd lets you log in
var ApiCmd = &cobra.Command{
	Use:     "api",
	Short:   "API server",
	Long:    "API server to aggregate data about your stars.",
	Aliases: []string{"api", "api-server"},
	Example: fmt.Sprintf("  %s classify", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		// do stuff

	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/api.go", 
			"cmd.name": 			"ApiCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	RootCmd.AddCommand(ApiCmd)
}

package action

import (
	"context" 																				// go-core
	"fmt"																					// go-core
	"github.com/roscopecoltran/sniperkit-limo/config" 										// app-config
	"github.com/roscopecoltran/sniperkit-limo/service" 										// svc-registry
	"github.com/spf13/cobra" 																// cli-cmd
	"github.com/sirupsen/logrus" 															// logs-logrus
	//"github.com/davecgh/go-spew/spew" 													// debug-print
	//"github.com/k0kubun/pp" 																// debug-print
)

// LoginCmd lets you log in
var LoginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Log in to a service",
	Long:    "Log in to the service specified by [--service] (default: github).",
	Example: fmt.Sprintf("  %s login", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		svc, err := getService() 															// Get the specified service and log in
		fatalOnError(err)

		token, err := svc.Login(ctx) 														// Save login access_token
		fatalOnError(err)

		config, err := getConfiguration() 													// Update configuration with token
		fatalOnError(err)

		config.GetService(service.Name(svc)).Token = token 									// Reload the service with the new token
		fatalOnError(config.WriteConfig()) 													// Write the updated config file

	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/login.go", 
			"cmd.name": 			"LoginCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	RootCmd.AddCommand(LoginCmd)
}

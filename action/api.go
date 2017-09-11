package actions

import (
	//"context"
	"fmt"
	"github.com/roscopecoltran/sniperkit-limo/config"
	"github.com/spf13/cobra"
	// bleve_http "github.com/blevesearch/bleve/http"
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
	RootCmd.AddCommand(ApiCmd)
}

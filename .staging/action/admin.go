package actions

import (
	"fmt"
	// "strconv"
	// "strings"
	// "github.com/blevesearch/bleve"
	"github.com/roscopecoltran/sniperkit-limo/config"
	// "github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/spf13/cobra"

)

// AdminCmd does a full-text gateway
var AdminCmd = &cobra.Command{
	Use:     "admin",
	Aliases: []string{"webui"},
	Short:   "Admin web-ui to manage your stars",
	Long:    "Admin web-ui to manage your stars and its attributes.",
	Example: fmt.Sprintf("  %s admin", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		// do stuff

	},
}

func init() {
	RootCmd.AddCommand(AdminCmd)
}

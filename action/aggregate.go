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

// AggregateCmd does a full-text gateway
var AggregateCmd = &cobra.Command{
	Use:     "aggregate <vcs uri>",
	Aliases: []string{"collect"},
	Short:   "Aggregate info on stars",
	Long:    "Perform a full-text aggregate search on your stars",
	Example: fmt.Sprintf("  %s aggregate robust", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		// do stuff

	},
}

func init() {
	RootCmd.AddCommand(AggregateCmd)
}

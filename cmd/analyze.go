package cmd

import (
	"fmt"
	// "strconv"
	// "strings"
	// "github.com/blevesearch/bleve"
	"github.com/hoop33/limo/config"
	// "github.com/hoop33/limo/model"
	"github.com/spf13/cobra"
)

// AnalyzeCmd does a full-text gateway
var AnalyzeCmd = &cobra.Command{
	Use:     "analyze <vcs uri>",
	Aliases: []string{"analyze", "augmented", "a"},
	Short:   "Analyze info on stars",
	Long:    "Perform an extended analyze on your stars",
	Example: fmt.Sprintf("  %s analyze robust", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		// do stuff

	},
}

func init() {
	RootCmd.AddCommand(AnalyzeCmd)
}

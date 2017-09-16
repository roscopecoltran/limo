// Copyright © 2017 Makoto Ito
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"fmt"																							// go-core
	//"github.com/roscopecoltran/sniperkit-limo/config" 											// app-config
	//"github.com/roscopecoltran/sniperkit-limo/service" 											// svc-registry
	//"github.com/roscopecoltran/sniperkit-limo/model" 												// data-models
	"github.com/spf13/cobra" 																		// cli-cmd
	"github.com/sirupsen/logrus" 																	// logs-logrus
	sim "github.com/ynqa/word-embedding/similarity" 												// ai-word-embed-word2vec
	"github.com/ynqa/word-embedding/utils" 															// ai-word-embed
	//"github.com/davecgh/go-spew/spew" 															// debug-print
	//"github.com/k0kubun/pp" 																		// debug-print
)

var (
	rank            	int
	inputVectorFile 	string
)

const helpSimilarityInstructions = `

	1. Please remove any protocol prefix from the url (eg. Golang like remote URIs)
	   - example: 'https://github.com/tensorflow/tensorflow' -> 'github.com/tensorflow/tensorflow'
	   - notes: 
	   		- the system will 'automatically-remove' the protocol part, like http:// or https://, 
	          but please note that for batch processing it will increase the overall performances 
	          if the input is already prefixed correctly and consistently.

	2. 'Similarity scoring' works with the following 'VCS' providers:
	   - GitHub (website: https://github.com)
	   - GitLab (website: https://gitlab.com)
	   - BitBucket (website: https://bitbucket.org)

`

// SimilarityCmd is the command for calculation of similarity.
var SimilarityCmd = &cobra.Command{
	Use:     fmt.Sprintf(" %s sim -i 'FILENAME' 'VCS_REPO_URI' \n %s", config.ProgramName, helpSimilarityInstructions),
	Short:   "Estimate the similarity between a starred repository.",
	Long:    "Estimate the similarity between a starred repository and all your bucket of starred repositories from github. gitlab and bitbucket.",
	Example: fmt.Sprintf("  %s sim -i ./shared/data/models/text/github_starred_vectors.txt github.com/tensorflow/tensorflow", config.ProgramName),
	Run: func(cmd *cobra.Command, args []string) {

		if !inputVectorFileIsExist() {
			utils.Fatal(fmt.Errorf("InputFile %s is not existed", inputVectorFile))
		}

		if len(args) == 1 {
			describe(args[0])
		} else {
			utils.Fatal(errors.New("Input a single word"))
		}
	},
}

func init() {
	SimilarityCmd.Flags().IntVarP(&rank, "rank", "r", 10, "Set number of the similar word list displayed")
	SimilarityCmd.Flags().StringVarP(&inputVectorFile, "input", "i", "example/word_vectors.txt",
		"Input path of a file written words' vector with libsvm format")
}

func describe(w string) {
	if err := sim.Load(inputVectorFile); err != nil {
		utils.Fatal(err)
	}

	if err := sim.Describe(w, rank); err != nil {
		utils.Fatal(err)
	}
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/similarity.go", 
			"cmd.name": 			"SimilarityCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	RootCmd.AddCommand(SimilarityCmd)
}

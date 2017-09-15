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
	"fmt"																								// go-core
	"github.com/roscopecoltran/sniperkit-limo/config" 													// app-config
	"github.com/roscopecoltran/sniperkit-limo/service" 													// svc-registry
	"github.com/roscopecoltran/sniperkit-limo/model" 													// data-models
	"github.com/spf13/cobra" 																			// cli-cmd
	"github.com/sirupsen/logrus" 																		// logs-logrus
	"github.com/ynqa/word-embedding/models/word2vec" 													// ai-word-embed-word2vec
	"github.com/ynqa/word-embedding/models/word2vec/model" 												// ai-word-embed-word2vec
	"github.com/ynqa/word-embedding/models/word2vec/opt" 												// ai-word-embed-word2vec
	"github.com/ynqa/word-embedding/utils"																// ai-word-embed
)

var (
	subModel, optimizer  string
	maxDepth, sampleSize int
	subsampleThreshold   float64
)

// Word2VecCmd is the command for word2vec.
var Word2VecCmd = &cobra.Command{
	Use:   "word2vec",
	Short: "Embed words using word2vec",
	Long:  "Embed words using word2vec",
	Example: fmt.Sprintf("  %s word2vec", config.ProgramName),
	//Aliases: []string{"w2c", "word2vec", "embed", "embed-word2vec"},
	Run: func(cmd *cobra.Command, args []string) {
		if !inputFileIsExist() {
			utils.Fatal(fmt.Errorf("InputFile %s is not existed", inputFile))
		} else if outputFileIsExist() {
			utils.Fatal(fmt.Errorf("OutputFile %s is already existed", outputFile))
		}

		if err := validateCommonParams(); err != nil {
			utils.Fatal(err)
		} else if err := validateWord2vecParams(); err != nil {
			utils.Fatal(err)
		}

		start()
	},
}

func init() {
	Word2VecCmd.Flags().AddFlagSet(GetCommonFlagSet())
	Word2VecCmd.Flags().StringVar(&subModel, "model", "cbow", "Set model from: skip-gram|cbow")
	Word2VecCmd.Flags().StringVar(&optimizer, "optimizer", "hs", "Set optimizer from: hs|ns")
	Word2VecCmd.Flags().IntVar(&maxDepth, "max-depth", 0, "Set number of times to track huffman tree, "+
		"max-depth=0 means tracking full path (using only hierarchical softmax)")
	Word2VecCmd.Flags().IntVar(&sampleSize, "negative", 5, "Set number of the samplings as negative instances "+
		"(using only negative sampling)")
	Word2VecCmd.Flags().Float64Var(&subsampleThreshold, "sample", 1.0e-3, "Set the threshold of subsampling")
}

// NewWord2Vec creates the word2vec struct.
func NewWord2Vec() word2vec.Word2Vec {
	return word2vec.Word2Vec{
		Common:             NewCommon(),
		Model:              NewModel(),
		Opt:                NewOptimizer(),
		SubSampleThreshold: subsampleThreshold,
	}
}

// NewOptimizer creates the optimizer of word2vec.
func NewOptimizer() (o word2vec.Optimizer) {
	switch optimizer {
	case "hs":
		o = opt.HierarchicalSoftmax{
			Common:   NewCommon(),
			MaxDepth: maxDepth,
		}
	case "ns":
		o = opt.NegativeSampling{
			Common:     NewCommon(),
			SampleSize: sampleSize,
		}
	}
	return
}

// NewModel creates the model of word2vec.
func NewModel() (m word2vec.Model) {
	switch subModel {
	case "skip-gram":
		m = model.SkipGram{
			Common: NewCommon(),
		}
	case "cbow":
		m = model.CBOW{
			Common: NewCommon(),
		}
	}
	return
}

func start() {
	w2v := NewWord2Vec()

	fmt.Print("Start PreTrain...\n")
	if err := w2v.PreTrain(); err != nil {
		utils.Fatal(err)
	}
	fmt.Print("Finish PreTrain\n")

	fmt.Printf("Vocabulary size: %d\n", word2vec.GetVocabulary())
	fmt.Printf("Number of words: %d\n", word2vec.GetWords())

	fmt.Print("Start Train...\n")
	if err := w2v.Run(); err != nil {
		utils.Fatal(err)
	}
	fmt.Print("Finish Train\n")

	if err := w2v.Save(); err != nil {
		utils.Fatal(err)
	}
	fmt.Print("Finish Save!\n")
}

func init() {
	log.WithFields(
		logrus.Fields{
			"src.file": 			"action/word2vec.go", 
			"cmd.name": 			"Word2VecCmd",
			"method.name": 			"init()", 
			"var.options": 			options, 
			}).Info("registering command...")
	RootCmd.AddCommand(Word2VecCmd)
}
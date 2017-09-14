package model

// https://github.com/tchajed/wordenc/tree/master/shortdict
// https://github.com/dustin/go-wikiparse
// 

type MachineLearningModelFiles struct {
	Gensim 		map[string]string
	Tensorflow  map[string]string
}

type MachineLearningModels struct {
	Word2Vec 	map[string]Word2VecModel
	Tensorflow  map[string]Tensorflow
}

type Tensorflow struct {
	words, size int
	vocab       []string
	M           []float32
}
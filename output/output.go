package output

import (
	"reflect"
	"strings"

	"github.com/hoop33/limo/config"
	"github.com/hoop33/limo/model"
)

// Output represents an output option
type Output interface {
	Configure(*config.OutputConfig)
	Inline(string)
	Info(string)
	Event(*model.Event)
	Error(string)
	Fatal(string)
	StarLine(*model.Star)
	Star(*model.Star)
	Repo(*model.Repo)
	Tag(*model.Tag)
	Topic(*model.Topic)
	Academic(*model.Academic)
	Software(*model.Software)
	Tree(*model.Tree)
	LanguageDetected(*model.LanguageDetected)
	Pkg(*model.Pkg)
	Readme(*model.Readme)
	Keyword(*model.Keyword)
	Pattern(*model.Pattern)
	Tick()
}

var outputs = make(map[string]Output)

func registerOutput(output Output) {
	parts := strings.Split(reflect.TypeOf(output).String(), ".")
	outputs[strings.ToLower(parts[len(parts)-1])] = output
}

// ForName returns the output for a given name
func ForName(name string) Output {
	if output, ok := outputs[name]; ok {
		return output
	}
	// We always want an output, so default to text
	return outputs["text"]
}

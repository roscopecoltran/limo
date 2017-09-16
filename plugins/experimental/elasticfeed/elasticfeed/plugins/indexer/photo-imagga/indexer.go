package photoimagga

import (
	"github.com/roscopecoltran/elasticfeed/common"
	"github.com/roscopecoltran/elasticfeed/workflow"
//	"github.com/roscopecoltran/elasticfeed/plugin/model"
)

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Indexer struct {
	config config
}

func (p *Indexer) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Indexer) Run(data interface{}) (interface{}, error) {
	return data, nil
}

func (p *Indexer) Cancel() {
}

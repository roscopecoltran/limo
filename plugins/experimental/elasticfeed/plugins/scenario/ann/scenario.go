package ann

import (
	"github.com/roscopecoltran/elasticfeed/common"
	"github.com/roscopecoltran/elasticfeed/workflow"
//	"github.com/roscopecoltran/elasticfeed/plugin/model"
)

type config struct {
	common.ElasticfeedConfig `mapstructure:",squash"`

	tpl *workflow.ConfigTemplate
}

type Scenario struct {
	config config
}

func (p *Scenario) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (p *Scenario) Run(data interface{}) (interface{}, error) {
	return data, nil
}

func (p *Scenario) Cancel() {
}

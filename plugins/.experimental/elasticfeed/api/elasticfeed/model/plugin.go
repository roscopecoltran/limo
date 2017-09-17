package model

import (
	pmodel "github.com/roscopecoltran/elasticfeed/plugin/model"
)

type PluginManager interface {

	LoadPipeline(name string) (pmodel.Pipeline, error)
}

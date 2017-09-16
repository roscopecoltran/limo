package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregateCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, AggregateCmd.Use)
}

func TestAggregateCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, AggregateCmd.Short)
}

func TestAggregateCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, AggregateCmd.Long)
}

func TestAggregateCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, AggregateCmd.Run)
}

package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyzeCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, AnalyzeCmd.Use)
}

func TestAnalyzeCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, AnalyzeCmd.Short)
}

func TestAnalyzeCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, AnalyzeCmd.Long)
}

func TestAnalyzeCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, AnalyzeCmd.Run)
}

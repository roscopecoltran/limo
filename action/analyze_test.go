package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

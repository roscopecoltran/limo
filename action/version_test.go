package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Use)
}

func TestVersionCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Short)
}

func TestVersionCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Long)
}

func TestVersionCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Run)
}

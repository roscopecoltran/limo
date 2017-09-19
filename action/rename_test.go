package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenameCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Use)
}

func TestRenameCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Short)
}

func TestRenameCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Long)
}

func TestRenameCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Run)
}

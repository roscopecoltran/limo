package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Use)
}

func TestLoginCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Short)
}

func TestLoginCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Long)
}

func TestLoginCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Run)
}

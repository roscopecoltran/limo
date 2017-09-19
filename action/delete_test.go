package action

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Use)
}

func TestDeleteCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Short)
}

func TestDeleteCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Long)
}

func TestDeleteCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Run)
}

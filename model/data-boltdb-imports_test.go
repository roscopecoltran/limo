package model

import (
	"io/ioutil"
	"os"
	"testing"
	// "github.com/roscopecoltran/sniperkit-limo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUp(t *testing.T) (*DatabaseDriver, func()) {
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err, "tmp file")

	filename := tmpFile.Name()
	boltClient := &DatabaseDriver{}
	err = boltClient.Open(filename)
	require.NoError(t, err, "open boltClient")

	return databaseDriver, func() {
		boltClient.Close()
		os.Remove(filename)
	}
}

func TestPaperRepository(t *testing.T) {
	databaseDriver, tearDown := setUp(t)
	defer tearDown()

	repo := NewPaperRepository(databaseDriver)

	// Not inserted yet -> id is 0
	id, err := repo.Get(1, "source 1", "ref 1")
	require.NoError(t, err, "get non inserted u1 s1 r1")
	assert.Equal(t, 0, id, "get non inserted u1 s1 r1 - id")

	err = repo.Save(1, 10, "source 1", "ref 1")
	require.NoError(t, err, "insert u1 p10 s1 r1")

	id, err = repo.Get(1, "source 1", "ref 1")
	require.NoError(t, err, "get u1 s1 r1")
	assert.Equal(t, 10, id, "get u1 s1 r1 - id")
}
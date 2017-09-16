package service

import (
	"github.com/disintegration/pak"
)

func mainChromiumPak() {
	// Read resources from file into memory
	p, err := pak.ReadFile("resources.pak")
	if err != nil {
		panic(err)
	}

	// Add or update resource
	p.Resources[12345] = []byte(`<!DOCTYPE html><html><h1>Hello from Go!</h1></html>`)

	// Write back to file
	err = pak.WriteFile("resources.pak", p)
	if err != nil {
		panic(err)
	}

}
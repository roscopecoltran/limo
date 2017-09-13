package main

import (
	"github.com/roscopecoltran/sniperkit-limo/service"
	"io"
	"os"
)

func main() {
	io.Copy(os.Stdout, service.NewReader(os.Stdin))
}


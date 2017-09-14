package main

import (
	"fmt"

	"github.com/jbowles/nlpt-cld2"
)

var s1 = "this sentence is in english dooode"

func this() {
}

func main() {
	res, err := cld2.SimpleDetect(s1)
	fmt.Printf("%v %v \n", res, err)
}
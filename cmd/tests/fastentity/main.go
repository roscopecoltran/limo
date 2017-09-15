package main

import (
	"fmt"
	"github.com/sajari/fastentity"
)

func main() {
	str := []rune("日 本語. Jack was a Golang developer from sydney. San Francisco, USA... Or so they say.")

	// Create a store
	store := fastentity.New("locations", "jobTitles")

	// Add single entities
	store.Add("locations", []rune("San Francisco, USA"))
	store.Add("jobTitles", []rune("golang developer"))

	// Add multiple (note: You don't need to initialise each group, they will be auto created if they don't exist)
	store.Add("skills", []rune("本語"), []rune("golang")) 

	results := store.FindAll(str)

	for group, entities := range results {
		fmt.Printf("Group: %s \n", group)
		for _, entity := range entities {
			fmt.Printf("\t-> %s\n", string(entity))
		}
	}
	/* Prints
	Group: locations
		-> San Francisco, USA
	Group: jobTitles
		-> golang developer
	Group: skills
		-> 本語
	*/
}
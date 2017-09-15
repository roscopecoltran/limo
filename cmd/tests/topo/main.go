package main

import (
	"fmt"
	"os"
	"github.com/fako1024/topo"
)

// List of all simple strings (to be sorted)
var stringsToSort = []string{
	"A", "B", "C", "D", "E", "F", "G", "H",
}

// List of dependencies
var stringDependencies = []topo.Dependency{
	topo.Dependency{Child: "B", Parent: "A"},
	topo.Dependency{Child: "B", Parent: "C"},
	topo.Dependency{Child: "B", Parent: "D"},
	topo.Dependency{Child: "A", Parent: "E"},
	topo.Dependency{Child: "D", Parent: "C"},
}

func main() {
	// Getter function to convert original elements to a generic type
	getter := func(i int) topo.Type {
		return stringsToSort[i]
	}

	// Setter function to restore the original type of the data
	setter := func(i int, val topo.Type) {
		stringsToSort[i] = val.(string)
	}

	// Perform topological sort
	if err := topo.Sort(stringsToSort, stringDependencies, getter, setter); err != nil {
		fmt.Printf("Error performing topological sort on slice of strings: %s\n", err)
		os.Exit(1)
	}

	// Print resulting Slice in order
	fmt.Println("Sorted list of strings:", stringsToSort)
	fmt.Println("The following dependencies were taken into account:")
	for _, dep := range stringDependencies {
		fmt.Println(dep)
	}
}
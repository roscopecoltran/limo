package main

import "gopkg.in/d4l3k/messagediff.v1"

type someStruct struct {
    A, b int // `testdiff:"ignore"`
    C []int
}

func main() {
    a := someStruct{1, 2, []int{1}}
    b := someStruct{1, 3, []int{1, 2}}
    diff, equal := messagediff.PrettyDiff(a, b)
    /*
        diff =
        `added: .C[1] = 2
        modified: .b = 3`

        equal = false
    */
}
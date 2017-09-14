package main

import (
    "fmt"
    "github.com/huandu/skiplist"
)

func main() {
    // Create a skip list with int key.
    list := skiplist.New(skiplist.Int)

    // Add some values. Value can be anything.
    list.Set(12, "hello world")
    list.Set(34, 56)

    // Get element by index.
    elem := list.Get(34) // Value is stored in elem.Value.
    fmt.Println(elem.Value)
    next := elem.Next()  // Get next element.
    fmt.Println(next.Value)

    // Or get value directly just like a map
    val, ok := list.GetValue(34)
    fmt.Println(val, ok)

    // Remove an element by index.
    list.Remove(34)
}
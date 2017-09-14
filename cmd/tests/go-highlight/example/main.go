package main

import "github.com/d4l3k/go-highlight"

func main() {
  highlight.Highlight("go", `
    package main

    import "fmt"

    func main() {
      fmt.Println("Duck!")
    }
  `)
  /*
    <keyword>package</keyword> main

    <keyword>import</keyword> <string>"fmt"</string>

    <keyword>func</keyword> main() {
      fmt.Println(<string>"Duck!"</string>)
    }
  */
}
package main

import (
  "fmt"
  "vlang-go/lang"
)

func input(prompt string) string {
  var i string
  fmt.Print(prompt)
  fmt.Scanln(&i)
  return i
}

func main() {
  for {
    text := input("vlang> ")
    result, error := lang.Run("<stdin>", text)
    if error != nil {
      fmt.Println(error.AsString())
    } else {
      var res []string
			for _, token := range result {
				res = append(res, fmt.Sprint(token.String()))
			}
      fmt.Println(res)
    }
  }
}

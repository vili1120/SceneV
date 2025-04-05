package main

import (
	"bufio"
	"fmt"
	"os"
	"SceneV/lang"
)

func input(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}


func main() {
  for {
    text := input("vlang> ")
    result, err := lang.Run("<stdin>", text)
    if err != nil {
      fmt.Println(err.AsString())
    } else {
      fmt.Println(result.String())
    }
  }
}

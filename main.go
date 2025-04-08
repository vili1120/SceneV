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
	globalSymbolTable := lang.NewSymbolTable()
  globalSymbolTable.Set("null", lang.NewNumber(0))
  globalSymbolTable.Set("true", lang.NewNumber(1))
  globalSymbolTable.Set("false", lang.NewNumber(0))

	for {
		text := input("SceneV> ")
		result, err := lang.Run("<stdin>", text, globalSymbolTable)
		if err != nil {
			fmt.Println(err.AsString())
		} else {
			fmt.Println(result)
		}
	}
}


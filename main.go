package main

import (
	"SceneV/lang"
	"os"
	"fmt"

	"github.com/chzyer/readline"
)

func input(prompt string) string {
	rl, err := readline.New(prompt)
	if err != nil {
		os.Exit(0)
	}
	defer rl.Close()
	line, err := rl.Readline()
	if err != nil {
		os.Exit(0)
	}
  return line
}


func main() {
	globalSymbolTable := lang.NewSymbolTable(nil)
  globalSymbolTable.Set("null", lang.NewNumber(0))
  globalSymbolTable.Set("true", lang.NewNumber(1))
  globalSymbolTable.Set("false", lang.NewNumber(0))

	for {
		text := input("SceneV> ")
    if text != "" {
		  result, err := lang.Run("<stdin>", text, globalSymbolTable)
		  if err != nil {
		  	fmt.Println(err.AsString())
		  } else if result != nil {
		  	fmt.Println(result.(lang.Val).String())
		  }
    }
	}
}


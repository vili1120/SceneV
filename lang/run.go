package lang

func Run(fn, text string) (any, *Error){
  lexer := NewLexer(fn, text)
  tokens, err := lexer.MakeTokens()
  if err != nil {
    return nil, err
  }

  parser := NewParser(tokens)
  ast := parser.Parse()

  if ast.error != nil {
    return nil, ast.error
  }

  interpreter := &Interpreter{}
  context := Context{
    "<program>",
    nil,
    nil,
  }
  result := interpreter.Visit(ast.node, context)
  
  return result.value, result.error
}

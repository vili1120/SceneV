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
  result := interpreter.Visit(ast.node)
  
  return result, nil
}

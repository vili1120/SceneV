package lang

func Run(fn, text string) (Node, *Error){
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
  
  return ast.node, nil
}

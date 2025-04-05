package lang

func Run(fn, text string) ([]Token, *Error){
  lexer := NewLexer(fn, text)
  tokens, err := lexer.MakeTokens()
  
  return tokens, err
}

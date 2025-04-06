package lang

const DIGITS = "0123456789"

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const LETTERS_DIGITS = LETTERS + DIGITS

const (
  INT           = "INT"
  FLOAT         = "FLOAT"

  IDENTIFIER    = "IDENTIFIER"
  KEYWORD       = "KEYWORD"
  
  PLUS          = "PLUS"
  MINUS         = "MINUS"
  MUL           = "MUL"
  DIV           = "DIV"
  POW           = "POW"

  EQ            = "EQ"

  LPAREN        = "LPAREN"
  RPAREN        = "RPAREN"

  EOF           = "EOF"
)

var KEYWORDS = []string{
  "var",
}

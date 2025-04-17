package lang

const DIGITS = "0123456789"

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const LETTERS_DIGITS = LETTERS + DIGITS

const (
  INT           = "INT"
  FLOAT         = "FLOAT"

  IDENTIFIER    = "IDENTIFIER"
  KEYWORD       = "KEYWORD"
  STRING        = "STRING"
  
  PLUS          = "PLUS"
  MINUS         = "MINUS"
  MUL           = "MUL"
  DIV           = "DIV"
  POW           = "POW"

  EQ            = "EQ"
  EE            = "EE"
  NE            = "NE"
  LT            = "LT"
  GT            = "GT"
  LTE           = "LTE"
  GTE           = "GTE"

  LPAREN        = "LPAREN"
  RPAREN        = "RPAREN"
  LBRACE        = "LBRACE"
  RBRACE        = "RBRACE"

  ARROW         = "ARROW"

  COMMA         = "COMMA"

  EOF           = "EOF"
)

var KEYWORDS = []string{
  "var",

  "and",
  "or",
  "not",

  "if",
  "elif",
  "else",

  "fn",

  "while",
  "for",
  "in",
}

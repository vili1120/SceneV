package lang

import (
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

func NewLexer(fn, text string) *Lexer {
  l := &Lexer {
    fn: fn,
    text: text,
    pos: Position{idx: -1, ln: 0, col: -1, fn: fn, ftxt: text},
    current_char: "",
  }
  l.advance()
  return l
}

type Lexer struct {
  fn string
  text string
  pos Position
  current_char string
}

func (l *Lexer) advance() {
  l.pos.Advance(l.current_char)
  if l.pos.idx < utf8.RuneCountInString(l.text) {
    r, _ := utf8.DecodeRuneInString(l.text[l.pos.idx:])
    l.current_char = string(r)
  } else {
    l.current_char = ""
  }
}

func (l *Lexer) MakeTokens() ([]Token, *Error) {
  tokens := []Token{}

  for l.current_char != "" {
    if strings.Contains(" \t", l.current_char) {
      l.advance()
    } else if strings.Contains(DIGITS, l.current_char) {
      tokens = append(tokens, l.MakeNumbers())
    }else if l.current_char == "+" {
      tokens = append(tokens, Token{type_: PLUS})
      l.advance()
    } else if l.current_char == "-" {
      tokens = append(tokens, Token{type_: MINUS})
      l.advance()
    } else if l.current_char == "*" {
      tokens = append(tokens, Token{type_: MUL})
      l.advance()
    } else if l.current_char == "/" {
      tokens = append(tokens, Token{type_: DIV})
      l.advance()
    } else if l.current_char == "(" {
      tokens = append(tokens, Token{type_: LPAREN})
      l.advance()
    } else if l.current_char == ")" {
      tokens = append(tokens, Token{type_: RPAREN})
      l.advance()
    } else {
      pos_start := l.pos.Copy()
      char := l.current_char
      l.advance()
      return nil, IllegalCharError(pos_start, l.pos, "'"+char+"'") 
    }
  }

  return tokens, nil
}

func (l *Lexer) MakeNumbers() Token {
  num_str := ""
  dot_count := 0

  for l.current_char != "" && strings.Contains(DIGITS+".", l.current_char) {
    if l.current_char == "." {
      if dot_count == 1 {
        break
      }
      dot_count += 1
      num_str += "."
    } else {
      num_str += l.current_char
    }
    l.advance()
  }

  if dot_count == 0 {
    num, err := strconv.Atoi(num_str)
    if err != nil { log.Fatal(err) }
    return Token{type_: INT, value: num}
  } else {
    num, err := strconv.ParseFloat(num_str, 64)
    if err != nil { log.Fatal(err) }
    return Token{type_: FLOAT, value: num}
  }
}

package lang

import (
	"log"
	"strconv"
	"strings"
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
	if l.pos.idx < len(l.text) {
		char := rune(l.text[l.pos.idx])
		l.current_char = string(char)
	} else {
		l.current_char = ""
	}
}


func (l *Lexer) MakeTokens() ([]Token, *Error) {
  tokens := []Token{}

  for l.current_char != "" {
    if l.current_char == " " || l.current_char == "\t" {
      l.advance()
      continue
    } else if strings.Contains(DIGITS, l.current_char) {
      tokens = append(tokens, l.MakeNumbers())
    }else if l.current_char == "+" {
      tokens = append(tokens, NewToken(PLUS, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "-" {
      tokens = append(tokens, NewToken(MINUS, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "*" {
      tokens = append(tokens, l.MakePower())
    } else if l.current_char == "/" {
      tokens = append(tokens, NewToken(DIV, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "(" {
      tokens = append(tokens, NewToken(LPAREN, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == ")" {
      tokens = append(tokens, NewToken(RPAREN, nil, &l.pos, nil))
      l.advance()
    } else {
      pos_start := l.pos.Copy()
      char := l.current_char
      l.advance()
      return nil, IllegalCharError(pos_start, l.pos, "'"+char+"'") 
    }
  }

  tokens = append(tokens, NewToken(EOF, nil, &l.pos, nil))
  return tokens, nil
}

func (l *Lexer) MakePower() Token {
  pos_start := l.pos.Copy()
  l.advance()

  if l.current_char == "*" {
    l.advance()
    return NewToken(POW, nil, &pos_start, &l.pos)
  }

  return NewToken(MUL, nil, &pos_start, &l.pos)
}


func (l *Lexer) MakeNumbers() Token {
  num_str := ""
  dot_count := 0
  pos_start := l.pos.Copy()

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
    return NewToken(INT, num, &pos_start, &l.pos)
  } else {
    num, err := strconv.ParseFloat(num_str, 64)
    if err != nil { log.Fatal(err) }
    return NewToken(FLOAT, num, &pos_start, &l.pos)
  }
}

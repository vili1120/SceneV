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
    } else if strings.Contains(LETTERS, l.current_char) {
      tokens = append(tokens, l.MakeIdentifier())
    } else if l.current_char == string('"') {
      tokens = append(tokens, l.MakeString())
    } else if l.current_char == "+" {
      tokens = append(tokens, NewToken(PLUS, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "-" {
      tokens = append(tokens, l.MakeArrow())
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
    } else if l.current_char == "{" {
      tokens = append(tokens, NewToken(LBRACE, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "}" {
      tokens = append(tokens, NewToken(RBRACE, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "," {
      tokens = append(tokens, NewToken(COMMA, nil, &l.pos, nil))
      l.advance()
    } else if l.current_char == "!" {
      tok, err := l.MakeNE()
      if err != nil { return nil, err }
      tokens = append(tokens, *tok)
    } else if l.current_char == "=" {
      tokens = append(tokens, l.MakeEquals())
    } else if l.current_char == "<" {
      tokens = append(tokens, l.MakeLT())
    } else if l.current_char == ">" {
      tokens = append(tokens, l.MakeGT())
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

func (l *Lexer) MakeIdentifier() Token {
  id_str := ""
  pos_start := l.pos.Copy()

  for l.current_char != "" && strings.Contains(LETTERS_DIGITS+"_", l.current_char) {
    id_str += l.current_char
    l.advance()
  }

  var tok_type string
  if contains(KEYWORDS, id_str) {
    tok_type = KEYWORD
  } else {
    tok_type = IDENTIFIER
  }
  return NewToken(tok_type, id_str, &pos_start, &l.pos)
}

func (l *Lexer) MakeNE() (*Token, *Error) {
  pos_start := l.pos.Copy()
  l.advance()

  if l.current_char == "=" {
    l.advance()
    tok := NewToken(NE, nil, &pos_start, &l.pos)
    return &tok, nil
  } else {
    l.advance()
    return nil, ExpectedCharError(pos_start, l.pos, "'=' (after '!')")
  }
}

func (l *Lexer) MakeEquals() Token {
  tok_type := EQ
  pos_start := l.pos.Copy()
  l.advance()

  if l.current_char == "=" {
    l.advance()
    tok_type = EE
  }

  return NewToken(tok_type, nil, &pos_start, &l.pos)
}

func (l *Lexer) MakeLT() Token {
  tok_type := LT
  pos_start := l.pos.Copy()
  l.advance()

  if l.current_char == "=" {
    l.advance()
    tok_type = LTE
  }

  return NewToken(tok_type, nil, &pos_start, &l.pos)
}

func (l *Lexer) MakeGT() Token {
  tok_type := GT
  pos_start := l.pos.Copy()
  l.advance()

  if l.current_char == "=" {
    l.advance()
    tok_type = GTE
  }

  return NewToken(tok_type, nil, &pos_start, &l.pos)
}

func (l *Lexer) MakeArrow() Token {
  tok_type := MINUS
  pos_start := l.pos.Copy()
  l.advance()

  if l.current_char == ">" {
    l.advance()
    tok_type = ARROW
  }

  return NewToken(tok_type, nil, &pos_start, &l.pos)
}

func (l *Lexer) MakeString() Token {
  str := ""
  posStart := l.pos.Copy()
  escapeChar := false
  l.advance()

  EscapeChars := map[string]string{
    "n": "\n",
    "t": "\t",
  }

  for l.current_char != "" && (l.current_char != string('"') || escapeChar) {
    if escapeChar {
      str += EscapeChars[l.current_char]
    } else {
      if l.current_char == "\\" {
        escapeChar = true
      } else {
        str += l.current_char
      }
    }
    l.advance()
    escapeChar = false
  }
  l.advance()
  return NewToken(STRING, str, &posStart, &l.pos)
}

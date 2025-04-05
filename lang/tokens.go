package lang

import "fmt"

func NewToken(type_ string, value any, pos_start, pos_end *Position) Token {
  var PosStart Position
  var PosEnd Position
  if pos_start != nil {
    PosStart = pos_start.Copy()
    PosEnd = pos_start.Copy()
    PosEnd.Advance()
  }
  if pos_end != nil {
    PosEnd = pos_end.Copy()
  }
  t := Token{
    type_: type_,
    value: value,
    PosStart: PosStart,
    PosEnd: PosEnd,
  }
  return t
}

type Token struct {
  type_ string
  value any
  PosStart Position
  PosEnd Position
}

func (t Token) String() string {
  if t.value != nil {
    return fmt.Sprintf("%v: %v", t.type_, t.value)
  } else {
    return fmt.Sprintf("%v", t.type_)
  }
}

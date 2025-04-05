package lang

import "fmt"

type Token struct {
  type_ string
  value any
}

func (t Token) String() string {
  if t.value != nil {
    return fmt.Sprintf("%v: %v", t.type_, t.value)
  } else {
    return fmt.Sprintf("%v", t.type_)
  }
}

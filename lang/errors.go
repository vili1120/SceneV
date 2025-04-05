package lang

import "fmt"

type Error struct {
  PosStart Position
  PosEnd Position
  ErrorName string
  Details string
}

func (e Error) AsString() string {
  result := fmt.Sprintf("%v: %v", e.ErrorName, e.Details)
  result += "\n"+fmt.Sprintf("(File: %v, Line: %v)", e.PosStart.fn, e.PosEnd.ln + 1)
  return result
}

func IllegalCharError(pos_start, pos_end Position, details string) *Error {
  return &Error{pos_start, pos_end, "Illegal Character", details}
}

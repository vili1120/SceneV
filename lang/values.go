package lang

import "fmt"

func NewNumber(value any) *Number {
  n := &Number{
    value: value,
  }
  n.SetPos(nil, nil)
  return n
}

type Number struct {
  value any
  PosStart, PosEnd *Position
}

func (n *Number) SetPos(pos_start, pos_end *Position) *Number {
  n.PosStart = pos_start
  n.PosEnd = pos_end
  return n
}

func (n *Number) Add(other *Number) *Number {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 + v2)
		case float64:
			return NewNumber(float64(v1) + v2)
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 + float64(v2))
		case float64:
			return NewNumber(v1 + v2)
		}
	}
	return nil
}

func (n *Number) Sub(other *Number) *Number {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 - v2)
		case float64:
			return NewNumber(float64(v1) - v2)
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 - float64(v2))
		case float64:
			return NewNumber(v1 - v2)
		}
	}
	return nil
}

func (n *Number) Mul(other *Number) *Number {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 * v2)
		case float64:
			return NewNumber(float64(v1) * v2)
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 * float64(v2))
		case float64:
			return NewNumber(v1 * v2)
		}
	}
	return nil
}

func (n *Number) Div(other *Number) *Number {
	switch v1 := n.value.(type) {
	case int:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(float64(v1) / float64(v2))
		case float64:
			return NewNumber(float64(v1) / v2)
		}
	case float64:
		switch v2 := other.value.(type) {
		case int:
			return NewNumber(v1 / float64(v2))
		case float64:
			return NewNumber(v1 / v2)
		}
	}
	return nil
}


func (n Number) String() string {
  return fmt.Sprintf("%v", n.value)
}

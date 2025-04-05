package lang

type Position struct {
  idx int
  ln int
  col int
  fn string
  ftxt string
}

func (p *Position) Advance(current_char string) {
  p.idx += 1
  p.col += 1

  if current_char == "\n" {
    p.ln += 1
    p.col = 0
  }
}

func (p Position) Copy() Position {
  return Position{
    p.idx,
    p.ln,
    p.col,
    p.fn,
    p.ftxt,
  }
}

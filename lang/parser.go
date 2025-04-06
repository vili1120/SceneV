package lang

func NewParser(tokens []Token) *Parser {
	p := &Parser{
		Tokens: tokens,
		TokIdx: -1,
	}
	p.advance()
	return p
}

type ParseResult struct {
	error *Error
	node  Node
  advance_count int
}

func (pr *ParseResult) register_advancement() {
  pr.advance_count += 1
}

func (pr *ParseResult) register(res *ParseResult) Node {
  pr.advance_count += res.advance_count
	if res.error != nil {
		pr.error = res.error
	}
	return res.node
}

func (pr *ParseResult) success(node Node) *ParseResult {
	pr.node = node
	return pr
}

func (pr *ParseResult) failure(err *Error) *ParseResult {
  if pr.error == nil || pr.advance_count == 0 {
	  pr.error = err
  }
	return pr
}

type Parser struct {
	Tokens     []Token
	TokIdx     int
	CurrentTok Token
}

func (p *Parser) advance() Token {
	p.TokIdx += 1
	if p.TokIdx < len(p.Tokens) {
		p.CurrentTok = p.Tokens[p.TokIdx]
	}
	return p.CurrentTok
}

func (p *Parser) Parse() *ParseResult {
	res := p.expr()
	if res.error == nil && p.CurrentTok.type_ != EOF {
		start := p.CurrentTok.PosStart
		end := p.CurrentTok.PosEnd
		for p.CurrentTok.type_ != EOF {
			end = p.CurrentTok.PosEnd
      res.register_advancement()
		  p.advance()
		}
		return res.failure(InvalidSyntaxError(
			start, end,
			"Expected '+', '-', '*', or '/'",
		))
	}
	return res
}

func (p *Parser) power() *ParseResult {
  return p.binOp(p.atom, []string{POW, POW}, p.factor)
}

func (p *Parser) atom() *ParseResult {
	res := &ParseResult{}
	tok := p.CurrentTok

  if contains([]string{INT, FLOAT}, tok.type_) {
    res.register_advancement()
		p.advance()
		return res.success(&NumberNode{Tok: tok, PosStart: tok.PosStart, PosEnd: tok.PosEnd})
	} else if tok.type_ == IDENTIFIER {
    res.register_advancement()
		p.advance()
    return res.success(&VarAccessNode{tok, tok.PosStart, tok.PosEnd})
  } else if tok.type_ == LPAREN {
    res.register_advancement()
		p.advance()
		expr := res.register(p.expr())
		if res.error != nil {
			return res
		}
		if p.CurrentTok.type_ != RPAREN {
			return res.failure(InvalidSyntaxError(
				p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
				"Expected ')'",
			))
		} else {
      res.register_advancement()
		  p.advance()
			return res.success(expr)
		}
	}

  return res.failure(InvalidSyntaxError(
    p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
    "Expected int, float, identifier, '+', '-', '('",
  ))
}

func (p *Parser) factor() *ParseResult {
	res := &ParseResult{}
	tok := p.CurrentTok

	if contains([]string{PLUS, MINUS}, tok.type_) {
    res.register_advancement()
		p.advance()
		factor := res.register(p.factor())
		if res.error != nil {
			return res
		}
		return res.success(&UnaryOpNode{tok, factor, tok.PosStart, getEndPos(factor)})
	}

	return p.power()
}

func (p *Parser) term() *ParseResult {
	return p.binOp(p.factor, []string{MUL, DIV}, nil)
}

func (p *Parser) expr() *ParseResult {
  res := ParseResult{}

  if p.CurrentTok.Matches(KEYWORD, "var") {
    res.register_advancement()
		p.advance()

    if p.CurrentTok.type_ != IDENTIFIER {
      return res.failure(InvalidSyntaxError(
        p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
        "Expected identifier",
      ))
    }
    var_name := p.CurrentTok
    res.register_advancement()
		p.advance()

    if p.CurrentTok.type_ != EQ {
      return res.failure(InvalidSyntaxError(
        p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
        "Expected '='",
      ))
    }
    res.register_advancement()
		p.advance()
    expr := res.register(p.expr())
    if res.error != nil { return &res }
    return res.success(&VarAssignNode{var_name, expr, var_name.PosStart, expr.GetPosEnd()})
  }

  node := res.register(p.binOp(p.term, []string{PLUS, MINUS}, nil))
  if res.error != nil {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'var', int, float, identifier, '+', '-', '('",
    ))
  }
  return res.success(node)
}

func (p *Parser) binOp(fna func() *ParseResult, ops []string, fnb func() *ParseResult) *ParseResult {
	if fnb == nil {
    fnb = fna
  }
  res := &ParseResult{}
	left := res.register(fna())
	if res.error != nil {
		return res
	}

	for contains(ops, p.CurrentTok.type_) {
		opTok := p.CurrentTok
    res.register_advancement()
		p.advance()
		right := res.register(fnb())
		if res.error != nil {
			return res
		}

		left = &BinOpNode{
			LeftNode:  left,
			OpTok:     opTok,
			RightNode: right,
			PosStart:  getStartPos(left),
			PosEnd:    getEndPos(right),
		}
	}

	return res.success(left)
}

func getStartPos(n Node) Position {
	return n.GetPosStart()
}

func getEndPos(n Node) Position {
	return n.GetPosEnd()
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

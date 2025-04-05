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
}

func (pr *ParseResult) register(res *ParseResult) Node {
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
	pr.error = err
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
      p.advance()
    }
		return res.failure(InvalidSyntaxError(
			start, end,
			"Expected '+', '-', '*', or '/'",
		))
	}
	return res
}

func (p *Parser) factor() *ParseResult {
	res := &ParseResult{}
	tok := p.CurrentTok
  
  if contains([]string{PLUS, MINUS}, tok.type_) {
    p.advance()
    factor := res.register(p.factor())
    if res.error != nil {
      return res
    }
    return res.success(UnaryOpNode{tok, factor})
  } else if contains([]string{INT, FLOAT}, tok.type_) {
		p.advance()
		return res.success(&NumberNode{Tok: tok})
	} else if tok.type_ == LPAREN {
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
      p.advance()
      return res.success(expr)
    }
  }

	return res.failure(InvalidSyntaxError(
		tok.PosStart, tok.PosEnd,
		"Expected int or float",
	))
}

func (p *Parser) term() *ParseResult {
	return p.binOp(p.factor, []string{MUL, DIV})
}

func (p *Parser) expr() *ParseResult {
	return p.binOp(p.term, []string{PLUS, MINUS})
}

func (p *Parser) binOp(fn func() *ParseResult, ops []string) *ParseResult {
  res := &ParseResult{}
  left := res.register(fn())
  if res.error != nil {
    return res
  }

  for contains(ops, p.CurrentTok.type_) {
    opTok := p.CurrentTok
    p.advance()
    right := res.register(fn())
    if res.error != nil {
      return res
    }

    left = &BinOpNode{
      LeftNode:  left,
      OpTok:     opTok,
      RightNode: right,
    }
  }

  return res.success(left)
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}


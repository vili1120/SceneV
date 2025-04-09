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
  return p.binOp(p.atom, []any{POW, POW}, p.factor)
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
	} else if tok.Matches(KEYWORD, "if") {
    if_expr := res.register(p.if_expr())
    if res.error != nil { return res }
    return res.success(if_expr)
  } else if tok.Matches(KEYWORD, "for") {
    for_expr := res.register(p.for_expr())
    if res.error != nil { return res }
    return res.success(for_expr)
  } else if tok.Matches(KEYWORD, "while") {
    while_expr := res.register(p.while_expr())
    if res.error != nil { return res }
    return res.success(while_expr)
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
	return p.binOp(p.factor, []any{MUL, DIV}, nil)
}

func (p *Parser) arith_expr() *ParseResult {
  return p.binOp(p.term, []any{PLUS, MINUS}, nil)
}

func (p *Parser) comp_expr() *ParseResult {
  res := ParseResult{}

  if p.CurrentTok.Matches(KEYWORD, "not") {
    op_tok := p.CurrentTok
    res.register_advancement()
    p.advance()

    node := res.register(p.comp_expr())
    if res.error != nil { return &res }
    return res.success(&UnaryOpNode{op_tok, node, op_tok.PosStart, op_tok.PosEnd})
  }

  node := res.register(p.binOp(p.arith_expr, []any{EE, NE, LT, GT, LTE, GTE}, nil))
  if res.error != nil {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'not', int, float, identifier, '+', '-', '('",
    ))
  }
  return res.success(node)
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

  node := res.register(p.binOp(p.comp_expr, []any{BinOpMatch{KEYWORD, "and"}, BinOpMatch{KEYWORD, "or"}}, nil))
  if res.error != nil {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'var', 'not', int, float, identifier, '+', '-', '('",
    ))
  }
  return res.success(node)
}

func (p *Parser) if_expr() *ParseResult {
  res := ParseResult{}
  cases := [][]Node{}
  var else_case Node

  if !p.CurrentTok.Matches(KEYWORD, "if") {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'if'",
    ))
  }
  res.register_advancement()
  p.advance()

  condition := res.register(p.expr())
  if res.error != nil { return &res }

  if p.CurrentTok.type_ != LBRACE {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected '{'",
    ))
  }
  res.register_advancement()
  p.advance()

  expr := res.register(p.expr())
  if res.error != nil { return &res }
  cases = append(cases, []Node{condition, expr})

  if p.CurrentTok.type_ != RBRACE {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected '}'",
    ))
  }
  res.register_advancement()
  p.advance()

  for p.CurrentTok.Matches(KEYWORD, "elif") {
    res.register_advancement()
    p.advance()

    condition := res.register(p.expr())
    if res.error != nil { return &res }

    if p.CurrentTok.type_ != LBRACE {
      return res.failure(InvalidSyntaxError(
        p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
        "Expected '{'",
      ))
    }
    res.register_advancement()
    p.advance()

    expr := res.register(p.expr())
    if res.error != nil { return &res }
    cases = append(cases, []Node{condition, expr})

    if p.CurrentTok.type_ != RBRACE {
      return res.failure(InvalidSyntaxError(
        p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
        "Expected '}'",
      ))
    }
    res.register_advancement()
    p.advance()
  }
  if p.CurrentTok.Matches(KEYWORD, "else") {
    res.register_advancement()
    p.advance()

    if p.CurrentTok.type_ != LBRACE {
      return res.failure(InvalidSyntaxError(
        p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
        "Expected '{'",
      ))
    }
    res.register_advancement()
    p.advance()

    else_case = res.register(p.expr())
    if res.error != nil { return &res }

    if p.CurrentTok.type_ != RBRACE {
      return res.failure(InvalidSyntaxError(
        p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
        "Expected '}'",
      ))
    }
    res.register_advancement()
    p.advance()
  }
  var pos_end Position
  if else_case != nil {
    pos_end = else_case.GetPosEnd()
  } else {
    pos_end = cases[len(cases) - 1][0].GetPosEnd()
  }
  return res.success(&IfNode{cases, else_case, cases[0][0].GetPosStart(), pos_end})
}

func (p *Parser) for_expr() *ParseResult {
  res := ParseResult{}
  
  if !p.CurrentTok.Matches(KEYWORD, "for") {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'for'",
    ))
  }
  res.register_advancement()
  p.advance()

  if p.CurrentTok.type_ != IDENTIFIER {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected identifier",
    ))
  }
  varName := p.CurrentTok
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

  startVal := res.register(p.expr())
  if res.error != nil { return &res }

  if !p.CurrentTok.Matches(KEYWORD, "in") {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'in'",
    ))
  }
  res.register_advancement()
  p.advance()

  endVal := res.register(p.expr())
  if res.error != nil { return &res }
  
  var stepVal Node
  if p.CurrentTok.type_ == ARROW {
    res.register_advancement()
    p.advance()

    stepVal = res.register(p.expr())
    if res.error != nil { return &res }
  } else {
    stepVal = nil
  }

  if p.CurrentTok.type_ != LBRACE {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected '{'",
    ))
  }
  res.register_advancement()
  p.advance()
  
  body := res.register(p.expr())
  if res.error != nil { return nil }

  if p.CurrentTok.type_ != RBRACE {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected '}'",
    ))
  }
  res.register_advancement()
  p.advance()

  posStart := varName.PosStart
  posEnd := body.GetPosEnd()
  return res.success(&ForNode{VarNameTok: varName, StartVal: startVal, EndVal: endVal, StepVal: stepVal, BodyNode: body, PosStart: posStart, PosEnd: posEnd})
}

func (p *Parser) while_expr() *ParseResult {
  res := ParseResult{}

  if !p.CurrentTok.Matches(KEYWORD, "while") {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected 'while'",
    ))
  }
  res.register_advancement()
  p.advance()

  condition := res.register(p.expr())
  if res.error != nil { return &res }

  if p.CurrentTok.type_ != LBRACE {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected '{'",
    ))
  }
  res.register_advancement()
  p.advance()

  body := res.register(p.expr())
  if res.error != nil { return &res }

  if p.CurrentTok.type_ != RBRACE {
    return res.failure(InvalidSyntaxError(
      p.CurrentTok.PosStart, p.CurrentTok.PosEnd,
      "Expected '}'",
    ))
  }
  res.register_advancement()
  p.advance()

  return res.success(&WhileNode{Cond: condition, BodyNode: body, PosStart: condition.GetPosStart(), PosEnd: body.GetPosEnd()})
}

func (p *Parser) binOp(fna func() *ParseResult, ops []any, fnb func() *ParseResult) *ParseResult {
	if fnb == nil {
    fnb = fna
  }
  res := &ParseResult{}
	left := res.register(fna())
	if res.error != nil {
		return res
	}

	for tokenMatchesFlexible(&p.CurrentTok, ops) {
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

type BinOpMatch struct {
	Type  string
	Value string
}

func tokenMatchesFlexible(tok *Token, ops []any) bool {
	for _, op := range ops {
		switch v := op.(type) {
		case string:
			if tok.type_ == v {
				return true
			}
		case BinOpMatch:
			if tok.Matches(v.Type, v.Value) {
				return true
			}
		default:
			panic("Invalid type in ops slice")
		}
	}
	return false
}

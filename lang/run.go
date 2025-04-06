package lang

func Run(fn string, text string, globalSymbolTable *SymbolTable) (any, *Error) {
	lexer := NewLexer(fn, text)
	tokens, err := lexer.MakeTokens()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	ast := parser.Parse()
	if ast.error != nil {
		return nil, ast.error
	}

	interpreter := &Interpreter{}
	context := Context{
		DisplayName: "<program>",
		SymbolTable: globalSymbolTable,
	}
	result := interpreter.Visit(ast.node, context)
	return result.value, result.error
}


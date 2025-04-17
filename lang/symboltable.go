package lang

type SymbolTable struct {
	Symbols map[string]Val
	Parent  *SymbolTable
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Symbols: map[string]Val{},
	}
}

func (st *SymbolTable) Get(name string) Val {
	value, ok := st.Symbols[name]
	if !ok && st.Parent != nil {
		return st.Parent.Get(name)
	}
	return value
}

func (st *SymbolTable) Set(name string, value Val) {
	st.Symbols[name] = value
}

func (st *SymbolTable) Remove(name string) {
  delete(st.Symbols, name)
}

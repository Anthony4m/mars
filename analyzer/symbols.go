package analyzer

import (
	"fmt"
	"mars/ast"
)

// Symbol represents a named entity in the program
type Symbol struct {
	Name       string
	Type       ast.Type
	IsMutable  bool
	IsFunction bool
	DeclaredAt ast.Node
	Scope      *Scope
}

// Scope represents a lexical scope in the program
type Scope struct {
	Parent  *Scope
	Symbols map[string]*Symbol
}

// SymbolTable manages scopes and symbols
type SymbolTable struct {
	CurrentScope *Scope
	GlobalScope  *Scope
}

// NewSymbolTable creates a new symbol table with a global scope
func NewSymbolTable() *SymbolTable {
	global := &Scope{
		Symbols: make(map[string]*Symbol),
	}
	return &SymbolTable{
		CurrentScope: global,
		GlobalScope:  global,
	}
}

// EnterScope creates a new scope and makes it current
func (st *SymbolTable) EnterScope() {
	st.CurrentScope = &Scope{
		Parent:  st.CurrentScope,
		Symbols: make(map[string]*Symbol),
	}
}

// ExitScope returns to the parent scope
func (st *SymbolTable) ExitScope() {
	if st.CurrentScope.Parent != nil {
		st.CurrentScope = st.CurrentScope.Parent
	}
}

// Define adds a new symbol to the current scope
func (st *SymbolTable) Define(name string, typ ast.Type, isMutable, isFunction bool, declaredAt ast.Node) error {
	if _, exists := st.CurrentScope.Symbols[name]; exists {
		return fmt.Errorf("symbol '%s' already defined in this scope", name)
	}

	st.CurrentScope.Symbols[name] = &Symbol{
		Name:       name,
		Type:       typ,
		IsMutable:  isMutable,
		IsFunction: isFunction,
		Scope:      st.CurrentScope,
		DeclaredAt: declaredAt,
	}
	return nil
}

// Resolve looks up a symbol by name, searching from current scope up to global
func (st *SymbolTable) Resolve(name string) (*Symbol, error) {
	scope := st.CurrentScope
	for scope != nil {
		if symbol, exists := scope.Symbols[name]; exists {
			return symbol, nil
		}
		scope = scope.Parent
	}
	return nil, fmt.Errorf("undefined symbol '%s'", name)
}

// IsGlobal returns true if the current scope is the global scope
func (st *SymbolTable) IsGlobal() bool {
	return st.CurrentScope == st.GlobalScope
}

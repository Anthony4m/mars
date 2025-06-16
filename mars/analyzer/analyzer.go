package analyzer

import (
	"fmt"
	"mars/ast"
)

// Analyzer performs semantic analysis on the AST
type Analyzer struct {
	errors    []Error
	symbols   *SymbolTable
	types     *TypeChecker
	immutable *ImmutabilityChecker
}

// New creates a new analyzer instance
func New() *Analyzer {
	return &Analyzer{
		symbols:   NewSymbolTable(),
		types:     NewTypeChecker(),
		immutable: NewImmutabilityChecker(),
	}
}

// Analyze performs semantic analysis on the given AST
func (a *Analyzer) Analyze(node ast.Node) error {
	// First pass: collect all declarations
	if err := a.collectDeclarations(node); err != nil {
		return err
	}

	// Second pass: type checking and immutability analysis
	if err := a.checkTypes(node); err != nil {
		return err
	}

	if err := a.checkImmutability(node); err != nil {
		return err
	}

	if len(a.errors) > 0 {
		return fmt.Errorf("semantic analysis failed with %d errors", len(a.errors))
	}

	return nil
}

// Error represents a semantic analysis error
type Error struct {
	Line   int
	Column int
	Msg    string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Line, e.Column, e.Msg)
}

// collectDeclarations performs the first pass of semantic analysis,
// collecting all declarations and building the symbol table
func (a *Analyzer) collectDeclarations(node ast.Node) error {
	// TODO: Implement declaration collection
	return nil
}

// checkTypes performs type checking on the AST
func (a *Analyzer) checkTypes(node ast.Node) error {
	// TODO: Implement type checking
	return nil
}

// checkImmutability verifies immutability rules
func (a *Analyzer) checkImmutability(node ast.Node) error {
	// TODO: Implement immutability checking
	return nil
}

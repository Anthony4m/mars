package ast

import (
	"testing"
)

func TestProgramString(t *testing.T) {
	program := &Program{
		Declarations: []Declaration{
			&VarDecl{
				IsMutable: true,
				Name:      &Identifier{Name: "x"},
				Type:      &Type{BaseType: "int"},
				Value:     &Literal{Type: "number", Value: 42},
			},
			&FuncDecl{
				Name: &Identifier{Name: "add"},
				Parameters: []*Parameter{
					{
						Name: &Identifier{Name: "a"},
						Type: &Type{BaseType: "int"},
					},
					{
						Name: &Identifier{Name: "b"},
						Type: &Type{BaseType: "int"},
					},
				},
				ReturnType: &Type{BaseType: "int"},
				Body: &BlockStatement{
					Statements: []Statement{
						&ReturnStatement{
							Value: &BinaryExpression{
								Left:     &Identifier{Name: "a"},
								Operator: "+",
								Right:    &Identifier{Name: "b"},
							},
						},
					},
				},
			},
			&StructDecl{
				Name: &Identifier{Name: "Point"},
				Fields: []*FieldDecl{
					{
						Name: &Identifier{Name: "x"},
						Type: &Type{BaseType: "int"},
					},
					{
						Name: &Identifier{Name: "y"},
						Type: &Type{BaseType: "int"},
					},
				},
			},
		},
	}

	if program.TokenLiteral() != "x" {
		t.Errorf("program.TokenLiteral wrong. got=%q", program.TokenLiteral())
	}
}

func TestVarDeclString(t *testing.T) {
	decl := &VarDecl{
		IsMutable: true,
		Name:      &Identifier{Name: "x"},
		Type:      &Type{BaseType: "int"},
		Value:     &Literal{Type: "number", Value: 42},
	}

	if decl.TokenLiteral() != "x" {
		t.Errorf("decl.TokenLiteral wrong. got=%q", decl.TokenLiteral())
	}
}

func TestFuncDeclString(t *testing.T) {
	decl := &FuncDecl{
		Name: &Identifier{Name: "add"},
		Parameters: []*Parameter{
			{
				Name: &Identifier{Name: "a"},
				Type: &Type{BaseType: "int"},
			},
			{
				Name: &Identifier{Name: "b"},
				Type: &Type{BaseType: "int"},
			},
		},
		ReturnType: &Type{BaseType: "int"},
		Body: &BlockStatement{
			Statements: []Statement{
				&ReturnStatement{
					Value: &BinaryExpression{
						Left:     &Identifier{Name: "a"},
						Operator: "+",
						Right:    &Identifier{Name: "b"},
					},
				},
			},
		},
	}

	if decl.TokenLiteral() != "add" {
		t.Errorf("decl.TokenLiteral wrong. got=%q", decl.TokenLiteral())
	}
}

func TestStructDeclString(t *testing.T) {
	decl := &StructDecl{
		Name: &Identifier{Name: "Point"},
		Fields: []*FieldDecl{
			{
				Name: &Identifier{Name: "x"},
				Type: &Type{BaseType: "int"},
			},
			{
				Name: &Identifier{Name: "y"},
				Type: &Type{BaseType: "int"},
			},
		},
	}

	if decl.TokenLiteral() != "Point" {
		t.Errorf("decl.TokenLiteral wrong. got=%q", decl.TokenLiteral())
	}
}

func TestBinaryExpressionString(t *testing.T) {
	expr := &BinaryExpression{
		Left:     &Identifier{Name: "a"},
		Operator: "+",
		Right:    &Identifier{Name: "b"},
	}

	if expr.TokenLiteral() != "a" {
		t.Errorf("expr.TokenLiteral wrong. got=%q", expr.TokenLiteral())
	}
}

func TestUnaryExpressionString(t *testing.T) {
	expr := &UnaryExpression{
		Operator: "-",
		Right:    &Identifier{Name: "x"},
	}

	if expr.TokenLiteral() != "-" {
		t.Errorf("expr.TokenLiteral wrong. got=%q", expr.TokenLiteral())
	}
}

func TestArrayLiteralString(t *testing.T) {
	expr := &ArrayLiteral{
		Elements: []Expression{
			&Literal{Type: "number", Value: 1},
			&Literal{Type: "number", Value: 2},
			&Literal{Type: "number", Value: 3},
		},
	}

	if expr.TokenLiteral() != "[" {
		t.Errorf("expr.TokenLiteral wrong. got=%q", expr.TokenLiteral())
	}
}

func TestStructLiteralString(t *testing.T) {
	expr := &StructLiteral{
		Type: &Identifier{Name: "Point"},
		Fields: []*FieldInit{
			{
				Name:  &Identifier{Name: "x"},
				Value: &Literal{Type: "number", Value: 1},
			},
			{
				Name:  &Identifier{Name: "y"},
				Value: &Literal{Type: "number", Value: 2},
			},
		},
	}

	if expr.TokenLiteral() != "Point" {
		t.Errorf("expr.TokenLiteral wrong. got=%q", expr.TokenLiteral())
	}
}

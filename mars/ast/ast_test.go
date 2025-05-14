package ast

import (
	"testing"
)

func TestProgramString(t *testing.T) {
	program := &Program{
		Declarations: []Declaration{
			&VarDecl{
				Mutable: true,
				Name:    &Identifier{Name: "x"},
				Type:    &Type{BaseType: "int"},
				Value:   &Literal{Token: "42", Value: 42},
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
		Mutable: true,
		Name:    &Identifier{Name: "x"},
		Type:    &Type{BaseType: "int"},
		Value:   &Literal{Token: "42", Value: 42},
	}

	if decl.TokenLiteral() != "x" {
		t.Errorf("decl.TokenLiteral wrong. got=%q", decl.TokenLiteral())
	}
}

func TestAssignmentStatement(t *testing.T) {
	stmt := &AssignmentStatement{
		Name: &Identifier{Name: "x"},
		Value: &BinaryExpression{
			Left:     &Identifier{Name: "x"},
			Operator: "+",
			Right:    &Literal{Token: "5", Value: 5},
		},
	}

	if stmt.TokenLiteral() != "=" {
		t.Errorf("stmt.TokenLiteral wrong. got=%q", stmt.TokenLiteral())
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

	if expr.TokenLiteral() != "+" {
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
			&Literal{Token: "1", Value: 1},
			&Literal{Token: "2", Value: 2},
			&Literal{Token: "3", Value: 3},
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
				Value: &Literal{Token: "1", Value: 1},
			},
			{
				Name:  &Identifier{Name: "y"},
				Value: &Literal{Token: "2", Value: 2},
			},
		},
	}

	if expr.TokenLiteral() != "Point" {
		t.Errorf("expr.TokenLiteral wrong. got=%q", expr.TokenLiteral())
	}
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		name     string
		typ      *Type
		expected string
	}{
		{
			name:     "base type",
			typ:      &Type{BaseType: "int"},
			expected: "int",
		},
		{
			name: "array type",
			typ: &Type{
				ArrayType: &Type{BaseType: "int"},
			},
			expected: "[]int",
		},
		{
			name: "pointer type",
			typ: &Type{
				PointerType: &Type{BaseType: "int"},
			},
			expected: "*int",
		},
		{
			name: "struct type",
			typ: &Type{
				StructType: &Identifier{Name: "Point"},
			},
			expected: "Point",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.typ.TokenLiteral(); got != tt.expected {
				t.Errorf("Type.TokenLiteral() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMemberExpression(t *testing.T) {
	expr := &MemberExpression{
		Object:   &Identifier{Name: "point"},
		Property: &Identifier{Name: "x"},
	}

	if expr.TokenLiteral() != "point" {
		t.Errorf("expr.TokenLiteral wrong. got=%q", expr.TokenLiteral())
	}
}

func TestBreakStatement(t *testing.T) {
	stmt := &BreakStatement{}

	if stmt.TokenLiteral() != "break" {
		t.Errorf("stmt.TokenLiteral wrong. got=%q", stmt.TokenLiteral())
	}
}

func TestContinueStatement(t *testing.T) {
	stmt := &ContinueStatement{}

	if stmt.TokenLiteral() != "continue" {
		t.Errorf("stmt.TokenLiteral wrong. got=%q", stmt.TokenLiteral())
	}
}

func TestPrintStatement(t *testing.T) {
	stmt := &PrintStatement{
		Expression: &Literal{Token: "42", Value: 42},
	}

	if stmt.TokenLiteral() != "log" {
		t.Errorf("stmt.TokenLiteral wrong. got=%q", stmt.TokenLiteral())
	}
}

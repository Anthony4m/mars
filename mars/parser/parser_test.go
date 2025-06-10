package parser

import (
	"fmt"
	"mars/ast"
	"mars/lexer"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain 1 statement. got=%d",
				len(program.Declarations))
		}

		stmt := program.Declarations[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Declarations) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Declarations))
	}
	stmt, ok := program.Declarations[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Declarations[0] is not ast.ExpressionStatement. got=%T",
			program.Declarations[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Name != "foobar" {
		t.Errorf("ident.Name not %s. got=%s", "foobar", ident.Name)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestVariableDeclarations(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedType       string
		expectedValue      interface{}
	}{
		{"x := 5;", "x", "", 5},                                     // type inference
		{"x : int = 5;", "x", "int", 5},                             // explicit type
		{"mut y := 10.5;", "y", "", 10.5},                           // mutable with inference
		{"mut y : float = 10.5;", "y", "float", 10.5},               // mutable with explicit type
		{"name := \"John\";", "name", "", "\"John\""},               // string with inference
		{"name : string = \"John\";", "name", "string", "\"John\""}, // string with explicit type
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain 1 declaration. got=%d",
				len(program.Declarations))
		}

		stmt := program.Declarations[0]
		if !testVariableDeclaration(t, stmt, tt.expectedIdentifier, tt.expectedType, tt.expectedValue) {
			return
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"x := 5;", 5.0},            // type inference
		{"x : int = 5;", 5.0},       // explicit type
		{"mut y := 42;", 42.0},      // mutable with inference
		{"mut y : int = 42;", 42.0}, // mutable with explicit type
		{"return 100;", 100.0},      // in return statement
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program has not enough declarations. got=%d",
				len(program.Declarations))
		}

		var expr ast.Expression
		switch stmt := program.Declarations[0].(type) {
		case *ast.VarDecl:
			expr = stmt.Value
		case *ast.ReturnStatement:
			expr = stmt.Value
		default:
			t.Fatalf("unexpected declaration type. got=%T", program.Declarations[0])
		}

		literal, ok := expr.(*ast.Literal)
		if !ok {
			t.Fatalf("expr not *ast.Literal. got=%T", expr)
		}
		if literal.Value != tt.expected {
			t.Errorf("literal.Value not %f. got=%f", tt.expected, literal.Value)
		}
		if literal.TokenLiteral() != fmt.Sprintf("%d", int(tt.expected)) {
			t.Errorf("literal.TokenLiteral not %s. got=%s",
				fmt.Sprintf("%d", int(tt.expected)), literal.TokenLiteral())
		}
	}
}

// TODO fix these errors
func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain %d statements. got=%d\n",
				1, len(program.Declarations))
		}

		stmt, ok := program.Declarations[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Declarations[0] is not ast.ExpressionStatement. got=%T",
				program.Declarations[0])
		}

		exp, ok := stmt.Expression.(*ast.UnaryExpression)
		if !ok {
			t.Fatalf("stmt is not ast.UnaryExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain %d statements. got=%d\n",
				1, len(program.Declarations))
		}

		stmt, ok := program.Declarations[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Declarations[0] is not ast.ExpressionStatement. got=%T",
				program.Declarations[0])
		}

		exp, ok := stmt.Expression.(*ast.BinaryExpression)
		if !ok {
			t.Fatalf("exp is not ast.BinaryExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.TokenLiteral()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Declarations) != 1 {
		t.Fatalf("program.Declarations does not contain %d statements. got=%d\n",
			1, len(program.Declarations))
	}

	stmt, ok := program.Declarations[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Declarations[0] is not ast.IfStatement. got=%T",
			program.Declarations[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("stmt.Alternative.Statements was not nil. got=%+v", stmt.Alternative)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func(x, y) { x + y; }`

	l := lexer.New(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Declarations) != 1 {
		t.Fatalf("program.Declarations does not contain %d statements. got=%d\n",
			1, len(program.Declarations))
	}

	stmt, ok := program.Declarations[0].(*ast.FuncDecl)
	if !ok {
		t.Fatalf("program.Declarations[0] is not ast.FuncDecl. got=%T",
			program.Declarations[0])
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(stmt.Parameters))
	}

	testLiteralExpression(t, stmt.Parameters[0].Name, "x")
	testLiteralExpression(t, stmt.Parameters[1].Name, "y")

	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(stmt.Body.Statements))
	}

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			stmt.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Declarations[0].(*ast.FuncDecl)

		if len(stmt.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(stmt.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, stmt.Parameters[i].Name, ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Declarations) != 1 {
		t.Fatalf("program.Declarations does not contain %d statements. got=%d\n",
			1, len(program.Declarations))
	}

	stmt, ok := program.Declarations[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Declarations[0])
	}

	exp, ok := stmt.Expression.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionCall. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStructDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			name   string
			fields []struct {
				name string
				typ  string
			}
		}
	}{
		{
			input: `struct Point {
				x : int;
				y : int;
			}`,
			expected: struct {
				name   string
				fields []struct {
					name string
					typ  string
				}
			}{
				name: "Point",
				fields: []struct {
					name string
					typ  string
				}{
					{name: "x", typ: "int"},
					{name: "y", typ: "int"},
				},
			},
		},
		{
			input: `struct Person {
				name : string;
				age : int;
				active : bool;
			}`,
			expected: struct {
				name   string
				fields []struct {
					name string
					typ  string
				}
			}{
				name: "Person",
				fields: []struct {
					name string
					typ  string
				}{
					{name: "name", typ: "string"},
					{name: "age", typ: "int"},
					{name: "active", typ: "bool"},
				},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain 1 statement. got=%d",
				len(program.Declarations))
		}

		stmt := program.Declarations[0]
		structDecl, ok := stmt.(*ast.StructDecl)
		if !ok {
			t.Fatalf("stmt not *ast.StructDecl. got=%T", stmt)
		}

		if structDecl.Name.Name != tt.expected.name {
			t.Errorf("structDecl.Name.Name not %s. got=%s",
				tt.expected.name, structDecl.Name.Name)
		}

		if len(structDecl.Fields) != len(tt.expected.fields) {
			t.Fatalf("structDecl.Fields length not %d. got=%d",
				len(tt.expected.fields), len(structDecl.Fields))
		}

		for i, field := range structDecl.Fields {
			expected := tt.expected.fields[i]
			if field.Name.Name != expected.name {
				t.Errorf("field.Name.Name not %s. got=%s",
					expected.name, field.Name.Name)
			}
			if field.Type.BaseType != expected.typ {
				t.Errorf("field.Type.BaseType not %s. got=%s",
					expected.typ, field.Type.BaseType)
			}
		}
	}
}

func TestStructLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			typeName string
			fields   []struct {
				name  string
				value interface{}
			}
		}
	}{
		{
			input: `Point{x: 5, y: 10}`,
			expected: struct {
				typeName string
				fields   []struct {
					name  string
					value interface{}
				}
			}{
				typeName: "Point",
				fields: []struct {
					name  string
					value interface{}
				}{
					{name: "x", value: 5},
					{name: "y", value: 10},
				},
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain 1 statement. got=%d",
				len(program.Declarations))
		}

		stmt := program.Declarations[0]
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", stmt)
		}

		structLit, ok := exprStmt.Expression.(*ast.StructLiteral)
		if !ok {
			t.Fatalf("expr not *ast.StructLiteral. got=%T", exprStmt.Expression)
		}

		if structLit.Type.Name != tt.expected.typeName {
			t.Errorf("structLit.Type.Name not %s. got=%s",
				tt.expected.typeName, structLit.Type.Name)
		}

		if len(structLit.Fields) != len(tt.expected.fields) {
			t.Fatalf("structLit.Fields length not %d. got=%d",
				len(tt.expected.fields), len(structLit.Fields))
		}

		for i, field := range structLit.Fields {
			expected := tt.expected.fields[i]
			if field.Name.Name != expected.name {
				t.Errorf("field.Name.Name not %s. got=%s",
					expected.name, field.Name.Name)
			}
			if !testLiteralExpression(t, field.Value, expected.value) {
				return
			}
		}
	}
}

func TestArrayTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected struct {
			size     *int
			elemType string
		}
	}{
		{
			input: `[5]int`,
			expected: struct {
				size     *int
				elemType string
			}{
				size:     func() *int { i := 5; return &i }(),
				elemType: "int",
			},
		},
		{
			input: `[]string`,
			expected: struct {
				size     *int
				elemType string
			}{
				size:     nil,
				elemType: "string",
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		typ := p.parseType()
		checkParserErrors(t, p)

		if typ == nil {
			t.Fatalf("parseType() returned nil")
		}

		if tt.expected.size == nil {
			if typ.ArraySize != nil {
				t.Errorf("typ.ArraySize not nil. got=%d", *typ.ArraySize)
			}
		} else {
			if typ.ArraySize == nil {
				t.Errorf("typ.ArraySize is nil")
			} else if *typ.ArraySize != *tt.expected.size {
				t.Errorf("typ.ArraySize not %d. got=%d",
					*tt.expected.size, *typ.ArraySize)
			}
		}

		if typ.ArrayType.BaseType != tt.expected.elemType {
			t.Errorf("typ.ArrayType.BaseType not %s. got=%s",
				tt.expected.elemType, typ.ArrayType.BaseType)
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{
			input:    `[1, 2, 3]`,
			expected: []interface{}{1, 2, 3},
		},
		{
			input:    `["hello", "world"]`,
			expected: []interface{}{"hello", "world"},
		},
		{
			input:    `[true, false]`,
			expected: []interface{}{true, false},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain 1 statement. got=%d",
				len(program.Declarations))
		}

		stmt := program.Declarations[0]
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", stmt)
		}

		arrayLit, ok := exprStmt.Expression.(*ast.ArrayLiteral)
		if !ok {
			t.Fatalf("expr not *ast.ArrayLiteral. got=%T", exprStmt.Expression)
		}

		if len(arrayLit.Elements) != len(tt.expected) {
			t.Fatalf("arrayLit.Elements length not %d. got=%d",
				len(tt.expected), len(arrayLit.Elements))
		}

		for i, elem := range arrayLit.Elements {
			if !testLiteralExpression(t, elem, tt.expected[i]) {
				return
			}
		}
	}
}

func TestPointerType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`*int`, "int"},
		{`*string`, "string"},
		{`*bool`, "bool"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		typ := p.parseType()
		checkParserErrors(t, p)

		if typ == nil {
			t.Fatalf("parseType() returned nil")
		}

		if typ.PointerType == nil {
			t.Fatalf("typ.PointerType is nil")
		}

		if typ.PointerType.BaseType != tt.expected {
			t.Errorf("typ.PointerType.BaseType not %s. got=%s",
				tt.expected, typ.PointerType.BaseType)
		}
	}
}

func TestUnsafeBlock(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: `unsafe {
				var x : *int;
				x = alloc(5);
				free(x);
			}`,
			expected: []string{
				"var x : *int;",
				"x = alloc(5);",
				"free(x);",
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Declarations) != 1 {
			t.Fatalf("program.Declarations does not contain 1 statement. got=%d",
				len(program.Declarations))
		}

		stmt := program.Declarations[0]
		unsafeBlock, ok := stmt.(*ast.UnsafeBlock)
		if !ok {
			t.Fatalf("stmt not *ast.UnsafeBlock. got=%T", stmt)
		}

		if len(unsafeBlock.Body.Statements) != len(tt.expected) {
			t.Fatalf("unsafeBlock.Body.Statements length not %d. got=%d",
				len(tt.expected), len(unsafeBlock.Body.Statements))
		}

		for i, stmt := range unsafeBlock.Body.Statements {
			if stmt.String() != tt.expected[i] {
				t.Errorf("stmt.String() not %s. got=%s",
					tt.expected[i], stmt.String())
			}
		}
	}
}

func TestParserErrors(t *testing.T) {
	tests := []struct {
		input          string
		expectedErrors []string
	}{
		{
			input: `struct Point {
				x : int
				y : int;
			}`,
			expectedErrors: []string{
				"expected semicolon after field declaration",
			},
		},
		{
			input: `func add(a : int, b : int) int {
				return a + b;
			}`,
			expectedErrors: []string{
				"expected '->' or '{' after function signature",
			},
		},
		{
			input: `var x : int = 5`,
			expectedErrors: []string{
				"expected semicolon after variable declaration",
			},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := NewParser(l)
		p.ParseProgram()

		if len(p.errors) != len(tt.expectedErrors) {
			t.Errorf("parser has %d errors, want %d", len(p.errors), len(tt.expectedErrors))
			for _, err := range p.errors {
				t.Errorf("parser error: %q", err)
			}
			continue
		}

		for i, err := range p.errors {
			if err != tt.expectedErrors[i] {
				t.Errorf("parser error %d = %q, want %q",
					i, err, tt.expectedErrors[i])
			}
		}
	}
}

func testVariableDeclaration(t *testing.T, s ast.Declaration, name string, expectedType string, expectedValue interface{}) bool {
	varDecl, ok := s.(*ast.VarDecl)
	if !ok {
		t.Errorf("s not *ast.VarDecl. got=%T", s)
		return false
	}

	if varDecl.Name.Name != name {
		t.Errorf("varDecl.Name.Name not '%s'. got=%s", name, varDecl.Name.Name)
		return false
	}

	if expectedType != "" && varDecl.Type.BaseType != expectedType {
		t.Errorf("varDecl.Type.BaseType not '%s'. got=%s", expectedType, varDecl.Type.BaseType)
		return false
	}

	if !testLiteralExpression(t, varDecl.Value, expectedValue) {
		return false
	}

	return true
}

// Helper functions for testing

func checkParserErrors(t *testing.T, p *parser) {
	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))
	for _, msg := range p.errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.BinaryExpression)
	if !ok {
		t.Errorf("exp is not ast.BinaryExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, float64(v))
	case float64:
		return testIntegerLiteral(t, exp, v)
	case string:
		// Check if it's a string literal (starts with quote) or an identifier
		if len(v) > 0 && v[0] == '"' {
			return testStringLiteral(t, exp, v[1:len(v)-1]) // Remove quotes
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value float64) bool {
	integ, ok := il.(*ast.Literal)
	if !ok {
		t.Errorf("il not *ast.Literal. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %f. got=%f", value, integ.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Name != value {
		t.Errorf("ident.Name not %s. got=%s", value, ident.Name)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("exp not *ast.Literal. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.Literal)
	if !ok {
		t.Errorf("exp not *ast.Literal. got=%T", exp)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value not %s. got=%s", value, str.Value)
		return false
	}

	return true
}

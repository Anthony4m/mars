package main

import (
	"fmt"
	"mars/ast"
	"mars/lexer"
	"mars/parser"
	"os"
	"path/filepath"
	"strings"
)

func formatFile(filename string) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", filename)
		os.Exit(1)
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	// Check file extension
	if filepath.Ext(filename) != ".mars" {
		fmt.Printf("Warning: File '%s' doesn't have .mars extension\n", filename)
	}

	// Lexical analysis
	l := lexer.New(string(content))

	// Parsing
	p := parser.NewParser(l)
	program := p.ParseProgram()

	// Check for parser errors
	errors := p.GetErrors()
	if errors != nil && errors.HasErrors() {
		fmt.Printf("Parse errors in '%s':\n", filename)
		for _, err := range errors.Errors() {
			fmt.Printf("  %s\n", err)
		}
		os.Exit(1)
	}

	// Format the program
	formatted := formatProgram(program)

	// Write back to file
	err = os.WriteFile(filename, []byte(formatted), 0644)
	if err != nil {
		fmt.Printf("Error writing formatted file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("Formatted '%s'\n", filename)
}

func formatProgram(program *ast.Program) string {
	var result strings.Builder

	for i, decl := range program.Declarations {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(formatDeclaration(decl, 0))
	}

	return result.String()
}

func formatDeclaration(decl ast.Declaration, indent int) string {
	switch d := decl.(type) {
	case *ast.VarDecl:
		return formatVarDecl(d, indent)
	case *ast.FuncDecl:
		return formatFuncDecl(d, indent)
	case *ast.StructDecl:
		return formatStructDecl(d, indent)
	case *ast.UnsafeBlock:
		return formatUnsafeBlock(d, indent)
	case *ast.BlockStatement:
		return formatBlockStatement(d, indent)
	case *ast.ExpressionStatement:
		return formatExpressionStatement(d, indent)
	default:
		return fmt.Sprintf("// Unknown declaration type: %T\n", decl)
	}
}

func formatVarDecl(vd *ast.VarDecl, indent int) string {
	var result strings.Builder

	// Add indentation
	result.WriteString(strings.Repeat("    ", indent))

	// Mutable keyword
	if vd.Mutable {
		result.WriteString("mut ")
	}

	// Variable name
	result.WriteString(vd.Name.Name)

	// Type annotation
	if vd.Type != nil {
		result.WriteString(" : ")
		result.WriteString(formatType(vd.Type))
	}

	// Assignment
	if vd.Value != nil {
		result.WriteString(" = ")
		result.WriteString(formatExpression(vd.Value))
	}

	result.WriteString(";")
	return result.String()
}

func formatFuncDecl(fd *ast.FuncDecl, indent int) string {
	var result strings.Builder

	// Add indentation
	result.WriteString(strings.Repeat("    ", indent))

	// Function keyword and name
	result.WriteString("func ")
	result.WriteString(fd.Name.Name)

	// Parameters
	result.WriteString("(")
	for i, param := range fd.Signature.Parameters {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(param.Name.Name)
		result.WriteString(": ")
		result.WriteString(formatType(param.Type))
	}
	result.WriteString(")")

	// Return type
	if fd.Signature.ReturnType != nil {
		result.WriteString(" -> ")
		result.WriteString(formatType(fd.Signature.ReturnType))
	}

	// Function body
	result.WriteString(" {\n")
	result.WriteString(formatBlockStatement(fd.Body, indent+1))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("}")

	return result.String()
}

func formatStructDecl(sd *ast.StructDecl, indent int) string {
	var result strings.Builder

	// Add indentation
	result.WriteString(strings.Repeat("    ", indent))

	// Struct keyword and name
	result.WriteString("struct ")
	result.WriteString(sd.Name.Name)
	result.WriteString(" {\n")

	// Fields
	for _, field := range sd.Fields {
		result.WriteString(strings.Repeat("    ", indent+1))
		result.WriteString(field.Name.Name)
		result.WriteString(": ")
		result.WriteString(formatType(field.Type))
		result.WriteString(";\n")
	}

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("}")

	return result.String()
}

func formatUnsafeBlock(ub *ast.UnsafeBlock, indent int) string {
	var result strings.Builder

	// Add indentation
	result.WriteString(strings.Repeat("    ", indent))

	result.WriteString("unsafe {\n")
	result.WriteString(formatBlockStatement(ub.Body, indent+1))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("}")

	return result.String()
}

func formatBlockStatement(bs *ast.BlockStatement, indent int) string {
	var result strings.Builder

	for _, stmt := range bs.Statements {
		result.WriteString(formatStatement(stmt, indent))
		result.WriteString("\n")
	}

	// Remove trailing newline
	if len(bs.Statements) > 0 {
		resultStr := result.String()
		if strings.HasSuffix(resultStr, "\n") {
			resultStr = resultStr[:len(resultStr)-1]
		}
		return resultStr
	}

	return result.String()
}

func formatStatement(stmt ast.Statement, indent int) string {
	switch s := stmt.(type) {
	case *ast.VarDecl:
		return formatVarDecl(s, indent)
	case *ast.AssignmentStatement:
		return formatAssignmentStatement(s, indent)
	case *ast.IfStatement:
		return formatIfStatement(s, indent)
	case *ast.ForStatement:
		return formatForStatement(s, indent)
	case *ast.ReturnStatement:
		return formatReturnStatement(s, indent)
	case *ast.PrintStatement:
		return formatPrintStatement(s, indent)
	case *ast.BreakStatement:
		return formatBreakStatement(s, indent)
	case *ast.ContinueStatement:
		return formatContinueStatement(s, indent)
	case *ast.BlockStatement:
		return formatBlockStatement(s, indent)
	case *ast.ExpressionStatement:
		return formatExpressionStatement(s, indent)
	default:
		return fmt.Sprintf("// Unknown statement type: %T", stmt)
	}
}

func formatAssignmentStatement(as *ast.AssignmentStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString(as.Name.Name)
	result.WriteString(" = ")
	result.WriteString(formatExpression(as.Value))
	result.WriteString(";")

	return result.String()
}

func formatIfStatement(is *ast.IfStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("if ")
	result.WriteString(formatExpression(is.Condition))
	result.WriteString(" {\n")
	result.WriteString(formatBlockStatement(is.Consequence, indent+1))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("}")

	if is.Alternative != nil {
		result.WriteString(" else {\n")
		result.WriteString(formatBlockStatement(is.Alternative, indent+1))
		result.WriteString("\n")
		result.WriteString(strings.Repeat("    ", indent))
		result.WriteString("}")
	}

	return result.String()
}

func formatForStatement(fs *ast.ForStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("for ")

	if fs.Init != nil {
		result.WriteString(formatStatement(fs.Init, 0))
		result.WriteString(" ")
	}

	result.WriteString(formatExpression(fs.Condition))

	if fs.Post != nil {
		result.WriteString("; ")
		result.WriteString(formatStatement(fs.Post, 0))
	}

	result.WriteString(" {\n")
	result.WriteString(formatBlockStatement(fs.Body, indent+1))
	result.WriteString("\n")
	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("}")

	return result.String()
}

func formatReturnStatement(rs *ast.ReturnStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("return")

	if rs.Value != nil {
		result.WriteString(" ")
		result.WriteString(formatExpression(rs.Value))
	}

	result.WriteString(";")
	return result.String()
}

func formatPrintStatement(ps *ast.PrintStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("log(")

	if ps.Expression != nil {
		result.WriteString(formatExpression(ps.Expression))
	}

	result.WriteString(");")
	return result.String()
}

func formatBreakStatement(bs *ast.BreakStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("break;")
	return result.String()
}

func formatContinueStatement(cs *ast.ContinueStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString("continue;")
	return result.String()
}

func formatExpressionStatement(es *ast.ExpressionStatement, indent int) string {
	var result strings.Builder

	result.WriteString(strings.Repeat("    ", indent))
	result.WriteString(formatExpression(es.Expression))
	result.WriteString(";")
	return result.String()
}

func formatExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Literal:
		return formatLiteral(e)
	case *ast.Identifier:
		return e.Name
	case *ast.BinaryExpression:
		return formatBinaryExpression(e)
	case *ast.UnaryExpression:
		return formatUnaryExpression(e)
	case *ast.FunctionCall:
		return formatFunctionCall(e)
	case *ast.ArrayLiteral:
		return formatArrayLiteral(e)
	case *ast.StructLiteral:
		return formatStructLiteral(e)
	case *ast.IndexExpression:
		return formatIndexExpression(e)
	case *ast.MemberExpression:
		return formatMemberExpression(e)
	default:
		return fmt.Sprintf("// Unknown expression type: %T", expr)
	}
}

func formatLiteral(lit *ast.Literal) string {
	switch v := lit.Value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", lit.Value)
	}
}

func formatBinaryExpression(be *ast.BinaryExpression) string {
	return fmt.Sprintf("%s %s %s",
		formatExpression(be.Left),
		be.Operator,
		formatExpression(be.Right))
}

func formatUnaryExpression(ue *ast.UnaryExpression) string {
	return fmt.Sprintf("%s%s", ue.Operator, formatExpression(ue.Right))
}

func formatFunctionCall(fc *ast.FunctionCall) string {
	var result strings.Builder

	result.WriteString(formatExpression(fc.Function))
	result.WriteString("(")

	for i, arg := range fc.Arguments {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(formatExpression(arg))
	}

	result.WriteString(")")
	return result.String()
}

func formatArrayLiteral(al *ast.ArrayLiteral) string {
	var result strings.Builder

	result.WriteString("[")
	for i, elem := range al.Elements {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(formatExpression(elem))
	}
	result.WriteString("]")

	return result.String()
}

func formatStructLiteral(sl *ast.StructLiteral) string {
	var result strings.Builder

	result.WriteString(sl.Type.Name)
	result.WriteString("{")

	for i, field := range sl.Fields {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(field.Name.Name)
		result.WriteString(": ")
		result.WriteString(formatExpression(field.Value))
	}

	result.WriteString("}")
	return result.String()
}

func formatIndexExpression(ie *ast.IndexExpression) string {
	return fmt.Sprintf("%s[%s]",
		formatExpression(ie.Object),
		formatExpression(ie.Index))
}

func formatMemberExpression(me *ast.MemberExpression) string {
	return fmt.Sprintf("%s.%s",
		formatExpression(me.Object),
		me.Property.Name)
}

func formatType(t *ast.Type) string {
	if t == nil {
		return ""
	}

	if t.BaseType != "" {
		return t.BaseType
	}

	if t.ArrayType != nil {
		if t.ArraySize != nil {
			return fmt.Sprintf("[%d]%s", *t.ArraySize, formatType(t.ArrayType))
		}
		return fmt.Sprintf("[]%s", formatType(t.ArrayType))
	}

	if t.PointerType != nil {
		return fmt.Sprintf("*%s", formatType(t.PointerType))
	}

	if t.StructName != "" {
		return t.StructName
	}

	return "unknown"
}

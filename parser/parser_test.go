package parser

import (
	"CricLang/ast"
	"CricLang/lexer"
	"testing"
)

func TestPlayerStatement(t *testing.T) {
	input := `player x = 5;
	player y = 10;
	player foobar = 838383;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testPlayerStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testPlayerStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "player" {
		t.Errorf("s.TokenLiteral not 'player', got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.PlayerStatement)
	if !ok {
		t.Errorf("s not *ast.PlayerStatement. got=%T", s)
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestSignalDecisionStatements(t *testing.T) {
	input := `
	signalDecision 5;
	signalDecision 10;
	signalDecision 1819;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn't contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		signalDecisionStmt, ok := stmt.(*ast.SignalDecisionStatement)
		if !ok {
			t.Errorf("stmt not *ast.SignalDecisionStatement. got=%T", stmt)
			continue
		}
		if signalDecisionStmt.TokenLiteral() != "signalDecision" {
			t.Errorf("signalDecisionStmt.TokenLiteral not 'signalDecision'. got=%q", signalDecisionStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program hasn't enough statements, got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not as.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

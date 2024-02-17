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

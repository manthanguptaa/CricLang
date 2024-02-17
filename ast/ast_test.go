package ast

import (
	"CricLang/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&PlayerStatement{
				Token: token.Token{Type: token.PLAYER, Literal: "player"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "player myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

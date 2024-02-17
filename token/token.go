package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "NO_BALL"
	EOF     = "MATCH_ENDED"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	//Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FIELD"
	PLAYER   = "PLAYER"
)

var keywords = map[string]TokenType{
	"field":  FUNCTION,
	"player": PLAYER,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

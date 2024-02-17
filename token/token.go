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
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION                = "FIELD"
	PLAYER                  = "PLAYER"
	TRUE                    = "NOT_OUT"
	FALSE                   = "OUT"
	APPEAL_IF               = "APPEAL_FOR_DECISION"
	APPEALOVERTURNED_ELSEIF = "APPEAL_OVERTURNED"
	APPEALREJECTED_ELSE     = "APPEAL_REJECTED"
	SIGNALDECISION_RETURN   = "SIGNAL_DECISION"
)

var keywords = map[string]TokenType{
	"field":            FUNCTION,
	"player":           PLAYER,
	"notOut":           TRUE,
	"out":              FALSE,
	"appeal":           APPEAL_IF,
	"appealOverturned": APPEALOVERTURNED_ELSEIF,
	"appealRejected":   APPEALREJECTED_ELSE,
	"signalDecision":   SIGNALDECISION_RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

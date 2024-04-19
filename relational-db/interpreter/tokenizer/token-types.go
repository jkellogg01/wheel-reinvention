package tokenizer

type Token struct {
	Content []byte
	Line    int
	Type    TokenType
}

type TokenType int

const (
	// one-character symbols
	TOKEN_SEMICOLON TokenType = iota
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_GREATER
	TOKEN_LESS
	TOKEN_SPLAT
	TOKEN_DOT
	TOKEN_COMMA
	TOKEN_EQUALS
	TOKEN_PLUS
	TOKEN_MINUS

	// two-character symbols
	TOKEN_GREATER_EQUALS
	TOKEN_LESS_EQUALS
	TOKEN_BANG_EQUALS

	// some special tokens
	TOKEN_EOF
	TOKEN_ERROR

	// keywords... almost definitely not going to implement all of these
	TOKEN_STRING
	TOKEN_IDENT
	TOKEN_NUMBER
	TOKEN_KEYWORD
)

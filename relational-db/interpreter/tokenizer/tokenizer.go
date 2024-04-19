package tokenizer

import "errors"

type Tokenizer struct {
	Data    []byte
	Start   int
	Current int
}

func NewTokenizer(data []byte) (Tokenizer, error) {
	if len(data) == 0 {
		return Tokenizer{}, errors.New("no data")
	}
	return Tokenizer{
		Data: data,
	}, nil
}

func (t *Tokenizer) Tokenize() (TokenQueue, error) {
	result := make(TokenQueue, 0)
	for t.Current < len(t.Data) {
		c := t.Data[t.Current]
		switch c {
		case '\'':
			result.Push(TOKEN_SQUOTE)
		case '"':
			result.Push(TOKEN_DQUOTE)
		case ';':
			result.Push(TOKEN_SEMICOLON)
		case '(':
			result.Push(TOKEN_LPAREN)
		case ')':
			result.Push(TOKEN_RPAREN)
		case '[':
			result.Push(TOKEN_LBRACK)
		case ']':
			result.Push(TOKEN_RBRACK)
		case '*':
			result.Push(TOKEN_SPLAT)
		case '.':
			result.Push(TOKEN_DOT)
		case ',':
			result.Push(TOKEN_COMMA)
		case '=':
			result.Push(TOKEN_EQUALS)
		case '+':
			result.Push(TOKEN_PLUS)
		case '-':
			result.Push(TOKEN_MINUS)
		}
	}
	return result, nil
}

type TokenQueue []Token

func (q *TokenQueue) Push(token Token) {
	*q = append(*q, token)
}

func (q *TokenQueue) Pop() (Token, error) {
	if len(*q) < 1 {
		return -1, errors.New("nothing to dequeue")
	}
	val := (*q)[0]
	*q = (*q)[1:]
	return val, nil
}

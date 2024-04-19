package tokenizer

import (
	"errors"
	"sync"
)

type Tokenizer struct {
	Data    []byte
	Start   int
	Current int
	Line    int
	Tokens  TokenList
}

func NewTokenizer(data []byte) (Tokenizer, error) {
	if len(data) == 0 {
		return Tokenizer{}, errors.New("no data")
	}
	return Tokenizer{
		Data:   data,
		Tokens: newTokenList(),
	}, nil
}

func (t *Tokenizer) Tokenize() error {
	for t.Current < len(t.Data) {
		c := t.Advance()
		if isNumber(c) {
			t.EmitLiteral(isNumber, TOKEN_NUMBER)
			continue
		}
		if isAlpha(c) {
			// we're just gonna decide it's a rule that identifiers need double quotes.
			// kinda sucks but makes my life easier.
			// may change later.
			t.EmitLiteral(isAlpha, TOKEN_KEYWORD)
			continue
		}
		switch c {
		case ' ', '\t':
			// skip whitespace
		case '\n':
			t.Line++
		case ';':
			t.Emit(TOKEN_SEMICOLON)
		case '(':
			t.Emit(TOKEN_LPAREN)
		case ')':
			t.Emit(TOKEN_RPAREN)
		case '*':
			t.Emit(TOKEN_SPLAT)
		case '.':
			t.Emit(TOKEN_DOT)
		case ',':
			t.Emit(TOKEN_COMMA)
		case '=':
			t.Emit(TOKEN_EQUALS)
		case '+':
			t.Emit(TOKEN_PLUS)
		case '-':
			t.Emit(TOKEN_MINUS)
		case '>':
			if t.Peek() == '=' {
				t.Emit(TOKEN_GREATER_EQUALS)
			} else {
				t.Emit(TOKEN_GREATER)
			}
		case '<':
			if t.Peek() == '=' {
				t.Emit(TOKEN_LESS_EQUALS)
			} else {
				t.Emit(TOKEN_LESS)
			}
		case '!':
			if t.Peek() == '=' {
				t.Emit(TOKEN_BANG_EQUALS)
			} else {
				return errors.New("invalid syntax")
			}
		case '\'':
			t.EmitLiteral(func(c byte) bool { return c != '\'' }, TOKEN_STRING)
		case '"':
			t.EmitLiteral(func(c byte) bool { return c != '"' }, TOKEN_IDENT)
		}
	}
	return nil
}

func (t *Tokenizer) EmitLiteral(bound func(byte) bool, ttype TokenType) error {
	for t.Current < len(t.Data) {
		c := t.Advance()
		if !bound(c) {
			t.Advance()
			t.Emit(ttype)
			return nil
		}
	}
	return errors.New("literal not terminated")
}

func (t *Tokenizer) Advance() byte {
	char := t.Data[t.Current]
	t.Current++
	return char
}

func (t *Tokenizer) Peek() byte {
	return t.Data[t.Current]
}

func (t *Tokenizer) Emit(ttype TokenType) {
	result := Token{
		Type:    ttype,
		Content: t.Data[t.Start:t.Current],
		Line:    t.Line,
	}
	t.Start = t.Current
	t.Tokens.Push(result)
}

type TokenList struct {
	Tokens []Token
	lock   sync.RWMutex
}

func newTokenList() TokenList {
	return TokenList{
		Tokens: make([]Token, 0),
		lock:   sync.RWMutex{},
	}
}

func (l *TokenList) Push(token Token) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Tokens = append(l.Tokens, token)
}

func (l *TokenList) Pop() (Token, error) {
	val, err := l.Peek()
	if err != nil {
		return val, err
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Tokens = l.Tokens[1:]
	return val, nil
}

func (l *TokenList) Peek() (Token, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if len(l.Tokens) < 1 {
		return Token{}, errors.New("nothing to return")
	}
	return l.Tokens[0], nil
}

func isNumber(c byte) bool {
	return '0' <= c && '9' >= c
}

func isAlpha(c byte) bool {
	return 'a' <= c && 'z' >= c ||
		'A' <= c && 'Z' >= c
}

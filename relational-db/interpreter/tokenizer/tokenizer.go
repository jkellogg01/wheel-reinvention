package tokenizer

import (
	"errors"
	"sync"

	"github.com/charmbracelet/log"
)

type Tokenizer struct {
	Data    []byte
	Start   int
	Current int
	Line    int
	tokens  TokenList
}

func NewTokenizer(data []byte) (Tokenizer, error) {
	if len(data) == 0 {
		return Tokenizer{}, errors.New("no data")
	}
	return Tokenizer{
		Data:   data,
		tokens: newTokenList(),
	}, nil
}

func (t *Tokenizer) Tokenize() error {
	for t.Current < len(t.Data) {
		c := t.Advance()
		log.Debug("Scanning...", "char", string(c), "idx", t.Current)
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
			t.EmitBound('\'', TOKEN_STRING)
		case '"':
			t.EmitBound('"', TOKEN_IDENT)
		}
	}
	return nil
}

func (t *Tokenizer) EmitBound(bound byte, ttype TokenType) error {
	for t.Current < len(t.Data) {
		c := t.Advance()
		if c == bound {
			t.Advance()
			t.Emit(ttype)
			return nil
		}
	}
	return errors.New("literal not terminated")
}

func (t *Tokenizer) EmitLiteral(bound func(c byte) bool, ttype TokenType) error {
	for t.Current < len(t.Data) {
		if !bound(t.Advance()) {
			t.Emit(ttype)
			return nil
		}
	}
	return errors.New("literal not terminated (this one shouln't happen)")
}

func (t *Tokenizer) Advance() byte {
	if t.Current >= len(t.Data) {
		log.Warn("tried to advance past data length")
		return 0
	}
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
	log.Debug("Emitting token", "type", ttype, "content", string(result.Content))
	t.Start = t.Current
	t.tokens.Push(result)
}

func (t *Tokenizer) GetTokens() []Token {
	t.tokens.lock.RLock()
	defer t.tokens.lock.RUnlock()
	return t.tokens.Tokens
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

package tokenizer

import "errors"

type Tokenizer struct {
	Data    []byte
	Start   int
	Current int
}

func (t *Tokenizer) Tokenize() (TokenQueue, error) {
	result := make(TokenQueue, 0)
	return result, nil
}

type TokenQueue []Token

func (s *TokenQueue) Enqueue(token Token) {
	*s = append(*s, token)
}

func (s *TokenQueue) Dequeue() (Token, error) {
	if len(*s) < 1 {
		return -1, errors.New("nothing to dequeue")
	}
	val := (*s)[0]
	*s = (*s)[1:]
	return val, nil
}

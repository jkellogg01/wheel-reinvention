package tokenizer

import (
	"testing"

	"github.com/charmbracelet/log"
)

func TestTokenizer(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	testPhrases := []string{
		"SELECT * FROM \"example_db\"",
	}

	phraseResults := [][]TokenType{
		{TOKEN_KEYWORD, TOKEN_SPLAT, TOKEN_KEYWORD, TOKEN_IDENT},
	}

	for i, phrase := range testPhrases {
		t.Logf("Testing phrase %d: %s", i, phrase)
		phrase := []byte(phrase)
		t.Log("Initializing tokenizer...")
		tknzr, err := NewTokenizer(phrase)
		if err != nil {
			t.Logf("Failed to initialize tokenizer: %s", err)
			t.FailNow()
		}
		t.Log("tokenizing...")
		err = tknzr.Tokenize()
		if err != nil {
			t.Logf("Tokenization failed: %s", err)
			t.Fail()
		}
		for j, token := range tknzr.GetTokens() {
			expect := phraseResults[i][j]
			t.Logf("Expected token of type %v, got %v", token.Type, expect)
			if token.Type != expect {
				t.Fail()
			}
		}
	}
}

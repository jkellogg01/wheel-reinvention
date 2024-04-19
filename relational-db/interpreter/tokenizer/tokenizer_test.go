package tokenizer

import "testing"

func TestTokenizer(t *testing.T) {
	testPhrases := []string{
		"SELECT * FROM 'example_db'",
	}

	phraseResults := [][]TokenType{
		{TOKEN_KEYWORD, TOKEN_SPLAT, TOKEN_KEYWORD, TOKEN_IDENT},
	}

	for i, phrase := range testPhrases {
		t.Logf("Testing phrase %d: %s", i, phrase)
		phrase := []byte(phrase)
		tknzr, err := NewTokenizer(phrase)
		if err != nil {
			t.Logf("Failed to initialize tokenizer: %s", err)
			t.FailNow()
		}
		err = tknzr.Tokenize()
		if err != nil {
			t.Logf("Tokenization failed: %s", err)
			t.Fail()
		}
		for j, token := range tknzr.GetTokens() {
			expect := phraseResults[i][j]
			if token.Type != expect {
				t.Logf("Expected token of type %v, instead got %v", token.Type, expect)
				t.Fail()
			}
		}
	}
}

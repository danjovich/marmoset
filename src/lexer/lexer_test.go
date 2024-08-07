package lexer

import (
	"marmoset/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `int five = 5;
int ten = 10;

int add(int x, int y) {
	x + y;
};

int result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
		return true;
} else {
		return false;
}

10 == 10;
10 != 9;
char[6] foobar = "foobar"
"foo bar"
"foo \" bar"
"foo \n\r\t bar"
[1, 2];
{"foo": "bar"}

void func() { }

bool var = true;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "int"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT_LIT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "int"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT_LIT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "int"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.INT, "int"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.INT, "int"},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.INT, "int"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT_LIT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT_LIT, "5"},
		{token.LT, "<"},
		{token.INT_LIT, "10"},
		{token.GT, ">"},
		{token.INT_LIT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT_LIT, "5"},
		{token.LT, "<"},
		{token.INT_LIT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT_LIT, "10"},
		{token.EQ, "=="},
		{token.INT_LIT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT_LIT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT_LIT, "9"},
		{token.SEMICOLON, ";"},
		{token.CHAR, "char"},
		{token.LBRACKET, "["},
		{token.INT_LIT, "6"},
		{token.RBRACKET, "]"},
		{token.IDENT, "foobar"},
		{token.ASSIGN, "="},
		{token.STRING_LIT, "foobar"},
		{token.STRING_LIT, "foo bar"},
		{token.STRING_LIT, "foo \" bar"},
		{token.STRING_LIT, "foo \n\r\t bar"},
		{token.LBRACKET, "["},
		{token.INT_LIT, "1"},
		{token.COMMA, ","},
		{token.INT_LIT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING_LIT, "foo"},
		{token.COLON, ":"},
		{token.STRING_LIT, "bar"},
		{token.RBRACE, "}"},
		{token.VOID, "void"},
		{token.IDENT, "func"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.BOOL, "bool"},
		{token.IDENT, "var"},
		{token.ASSIGN, "="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

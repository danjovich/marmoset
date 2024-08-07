package token

type TokenType string

// in a production environment it makes sense to attach filenames and line numbers to tokens,
// to better track down lexing and parsing errors. So it would be better to initialize the
// lexer with an io.Reader and the filename
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT      = "IDENT"   // add, foobar, x, y, ...
	INT_LIT    = "INT_LIT" // 1, 2, 3, -1, ...
	STRING_LIT = "STRING_LIT"
	CHAR_LIT   = "CHAR_LIT"

	// Types
	INT  = "INT"
	CHAR = "CHAR"
	VOID = "VOID"
	BOOL = "BOOL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	// TODO: Remove let and function
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	// types
	"int":  INT,
	"char": CHAR,
	"void": VOID,
	"bool": BOOL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

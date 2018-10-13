package lexer

// TokenType maps token to type
type TokenType int

// Token defines individual token
type Token struct {
	Type     TokenType
	Contents string
	//Where    Span
}

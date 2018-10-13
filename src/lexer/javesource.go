package lexer

// JaveSource defines the input jave file
type JaveSource struct {
	Path     string
	Name     string
	Contents []rune
	NewLines []int
	Tokens   []*Token
}

package parser

type ParseResult interface {
	Command() string
}

type Parser interface {
	Accept([]byte) bool
	Parse([]byte) (ParseResult, error)
}

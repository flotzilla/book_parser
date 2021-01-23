package parser

import (
	src "book_parser/src"
)

type XMLParserResultHandler struct {
}

func (handler *XMLParserResultHandler) Handle(result *src.ParseResult) bool {
	return false
}

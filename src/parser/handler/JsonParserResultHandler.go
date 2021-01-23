package parser_handler

import (
	src "book_parser/src"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type JsonParserResultHandler struct {
}

func (handler *JsonParserResultHandler) Handle(result *src.ParseResult) bool {
	logrus.Trace("Will marshal to json")
	marshaled, err := json.Marshal(result)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Trace(string(marshaled))

	return false
}

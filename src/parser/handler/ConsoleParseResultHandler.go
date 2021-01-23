package parser_handler

import (
	"book_parser/src"
	"github.com/sirupsen/logrus"
)

type ConsoleParseResultHandler struct {
}

func (handler *ConsoleParseResultHandler) Handle(pr *src.ParseResult) bool {
	if pr == nil {
		return false
	}
	logrus.Debug("\tScan results:")
	logrus.Debug("Machine Id: " + pr.MachineId)
	logrus.Debug("ParseId: " + pr.ParseId)
	logrus.Debug("Directory: " + pr.FilePath)
	logrus.Debug("Scan duration: " + pr.Duration.String())
	logrus.Debug("\tBooks: ")
	for _, el := range pr.Books {
		logrus.Debug(el.BookInfo.Title, el.BookFile.Name)
	}

	return true
}

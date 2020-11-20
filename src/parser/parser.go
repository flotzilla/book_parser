package parser

import (
	"book_parser/src"
	"book_parser/src/parser/pdf"
	"errors"
	"fmt"
)

const (
	extPDF = "pdf"
	extFB2 = "fb2"
)

var (
	errInvalidFileTypeParser = errors.New("cannot handle this file type")
)

type ParseResult struct {
	Books  []src.Book
	Errors []error
	// time start
	// time end
	// books counter ??
}

func Parse(scanResult *src.ScanResult) *ParseResult {

	fmt.Println("Starting to parse")
	pr := ParseResult{}

	if len(scanResult.Books) == 0 {
		return &pr
	}

	for _, bookFile := range scanResult.Books {
		book := src.Book{}

		// handle symlinks
		if bookFile.IsSymLink {
			continue
		}

		switch bookFile.Ext {
		case extPDF:
			book.BookInfo = *pdf.Parse(&bookFile)
		case extFB2:
			// handle fb2
		default:
			pr.Errors = append(pr.Errors, errInvalidFileTypeParser)
		}
		pr.Books = append(pr.Books, book)
	}

	return &pr
}

package parser

import (
	"book_parser/src"
	_ "book_parser/src/logging"
	"book_parser/src/parser/pdf"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	extPDF = "pdf"
	extFB2 = "fb2"
)

var (
	errInvalidFileTypeParser = errors.New("cannot handle this file type")
	errCannotParseFile       = errors.New("cannot parse book file")
	symlinkError             = errors.New("book is a symlink")
)

type ParseResult struct {
	Books  []src.Book
	Errors []error
	// time start
	// time end
	// books counter ??
}

func Parse(scanResult *src.ScanResult) *ParseResult {

	logrus.Debug("Starting to parse")
	pr := ParseResult{}

	if len(scanResult.Books) == 0 {
		return &pr
	}

	var wg sync.WaitGroup

	booksCount := len(scanResult.Books)
	wg.Add(booksCount)
	booksChan := make(chan src.Book, booksCount)
	errorsChan := make(chan error, booksCount)

	for _, bookFile := range scanResult.Books {
		go parseBook(&wg, bookFile, booksChan, errorsChan)
	}

	logrus.Info("Waiting for parsing workers to finish")
	wg.Wait()
	pr.Books = append(pr.Books, <-booksChan)
	pr.Errors = append(pr.Errors, <-errorsChan)
	logrus.Info("Parsing workers finished")

	return &pr
}
func parseBook(wg *sync.WaitGroup, bookFile src.BookFile, bookChan chan src.Book, errorChan chan error) {
	logrus.Debug("Worker for ", bookFile.FilePath, " started")
	var (
		bookInfo *src.BookInfo
		err      error
	)

	defer func() {
		logrus.Debug("Worker for ", bookFile.FilePath, " finished")
		wg.Done()
	}()

	if bookFile.IsSymLink {
		errorChan <- symlinkError // ++ filename
	}

	switch bookFile.Ext {
	case extPDF:
		bookInfo, err = pdf.Parse(&bookFile)
	case extFB2:
		// TODO  handle fb2
		logrus.Warn("fb2 parse in process")
	default:
		errorChan <- errInvalidFileTypeParser
	}

	if err != nil {
		errorChan <- err
		return
	}

	if bookInfo == nil {
		errorChan <- errCannotParseFile
		return
	}

	bookChan <- src.Book{
		BookFile: bookFile,
		BookInfo: *bookInfo,
	}
}

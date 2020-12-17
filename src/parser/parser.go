package parser

import (
	"book_parser/src"
	_ "book_parser/src/logging"
	"book_parser/src/parser/types"
	"book_parser/src/utils"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	extPDF = "pdf"
	extFB2 = "fb2"
)

var (
	errInvalidFileTypeParser = errors.New("cannot handle this file type")
	errCannotParseFile       = errors.New("cannot parse book file")
	symlinkError             = errors.New("book is a symlink")

	cnf *src.Config
)

func Parse(scanResult *src.ScanResult, config *src.Config) *src.ParseResult {

	logrus.Debug("Starting to parse")
	pr := src.ParseResult{}
	cnf = config

	if len(scanResult.Books) == 0 {
		return &pr
	}

	var wg sync.WaitGroup

	booksCount := len(scanResult.Books)
	logrus.Debug("Books count: ", booksCount)
	wg.Add(booksCount)
	booksChan := make(chan src.Book, booksCount)
	errorsChan := make(chan error, booksCount)

	startTime := time.Now()
	for _, bookFile := range scanResult.Books {
		go parseBook(&wg, bookFile, booksChan, errorsChan)
	}

	logrus.Info("Waiting for parsing workers to finish")
	wg.Wait()
	close(booksChan)
	close(errorsChan)

	elapsed := time.Since(startTime)

	for el := range booksChan {
		pr.Books = append(pr.Books, el)
	}
	for el := range errorsChan {
		pr.Errors = append(pr.Errors, el)
	}

	logrus.Info("Parsing workers finished. Elapsed time: ", elapsed)
	logrus.Info("ParsedBooks: ", len(pr.Books), ". Errors: ", len(pr.Errors))

	return &pr
}

func parseBook(wg *sync.WaitGroup, bookFile src.BookFile, bookChan chan src.Book, errorChan chan error) {
	logrus.Debug("Worker started for ", bookFile.FilePath)
	var (
		bookInfo *src.BookInfo
		err      error
	)

	defer func() {
		logrus.Debug("Worker finished for ", bookFile.FilePath)
		wg.Done()
	}()

	if bookFile.IsSymLink {
		errorChan <- src.ParseError{PreviousError: symlinkError, FileName: bookFile.FilePath}
		return
	}

	if !utils.IsStringInSlice(bookFile.Ext, cnf.ScanExt) {
		logrus.Debug("Skipped ext,", cnf.ScanExt, " for file ", bookFile.FilePath)
		return
	}

	switch bookFile.Ext {
	case extPDF:
		bookInfo, err = types.Parse(&bookFile)
	case extFB2:
		// TODO  handle fb2
		logrus.Warn("fb2 parse in process")
	default:
		errorChan <- src.ParseError{PreviousError: errInvalidFileTypeParser, FileName: bookFile.FilePath}
		return
	}

	if err != nil {
		errorChan <- src.ParseError{PreviousError: err, FileName: bookFile.FilePath}
		return
	}

	if bookInfo == nil {
		errorChan <- src.ParseError{PreviousError: errCannotParseFile, FileName: bookFile.FilePath}
		return
	}

	bookChan <- src.Book{
		BookFile: bookFile,
		BookInfo: *bookInfo,
	}
}

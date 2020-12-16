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
	wg.Add(booksCount)
	booksChan := make(chan src.Book, booksCount)
	errorsChan := make(chan error, booksCount)

	startTime := time.Now()
	for _, bookFile := range scanResult.Books {
		go parseBook(&wg, bookFile, booksChan, errorsChan)
	}

	logrus.Info("Waiting for parsing workers to finish")
	wg.Wait()
	elapsed := time.Since(startTime)
	pr.Books = append(pr.Books, <-booksChan)
	pr.Errors = append(pr.Errors, <-errorsChan)
	logrus.Info("Parsing workers finished. Elapsed time: ", elapsed)

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
		errorChan <- symlinkError
		return
	}

	if !utils.IsStringInSlice(bookFile.Ext, cnf.ScanExt) {
		logrus.Debug("Skipped ext,", cnf.ScanExt, " for file ", bookFile.Name)
		//return
	}

	switch bookFile.Ext {
	case extPDF:
		bookInfo, err = types.Parse(&bookFile)
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

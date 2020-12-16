package parser

import (
	"book_parser/src"
	"book_parser/src/parser/pdf"
	"errors"
	"fmt"
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

	fmt.Println("Starting to parse")
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

	fmt.Println("Waiting for parsing workers to finish")
	wg.Wait()
	pr.Books = append(pr.Books, <-booksChan)
	pr.Errors = append(pr.Errors, <-errorsChan)
	fmt.Println("Parsing workers finished")

	return &pr
}
func parseBook(wg *sync.WaitGroup, bookFile src.BookFile, bookChan chan src.Book, errorChan chan error) {
	fmt.Println("Worker for ", bookFile.FilePath, " started")
	var (
		bookInfo *src.BookInfo
		err      error
	)

	defer func() {
		fmt.Println("Worker for ", bookFile.FilePath, " finished")
		if bookInfo != nil {
			fmt.Println(bookInfo.Title)
		}
		fmt.Println("=======")
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
		fmt.Println("fb2 parse in process")
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

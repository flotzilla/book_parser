package parser

import (
	"book_parser/src"
	cnf "book_parser/src/config"
	_ "book_parser/src/logging"
	fb2_parser "book_parser/src/parser/types/fb2"
	pdf_parser "book_parser/src/parser/types/pdf"
	"book_parser/src/utils"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/denisbrodbeck/machineid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	extPDF = "pdf"
	extFB2 = "fb2"
)

type Parser struct{}

var (
	errInvalidFileTypeParser = errors.New("cannot handle this file type")
	errCannotParseFile       = errors.New("cannot parse book file")
	symlinkError             = errors.New("book is a symlink")
	emptyScanBookCountError  = errors.New("scan result is empty")
	generateIdError          = errors.New("cannot generateMachineID")
)

func New() Parser {
	return Parser{}
}

// GenerateParseId
// Will create parse unique hash from this data:
// * mashine ID
// * start unix time
// * configuration hash
// * duration
// * count of find books
// * scan filepath
func (parser Parser) GenerateParseId(result *src.ParseResult, config *cnf.Config) string {
	hasher := sha256.New()
	hasher.Write([]byte(result.MachineId))

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(result.StartTime))
	hasher.Write(b)

	hasher.Write([]byte(config.GetConfigHash()))
	hasher.Write([]byte(result.Duration.String()))

	b1 := make([]byte, 8)
	binary.LittleEndian.PutUint32(b1, uint32(len(result.Books)))
	hasher.Write(b1)

	hasher.Write([]byte(result.FilePath))

	hash := hex.EncodeToString(hasher.Sum(nil))
	logrus.Trace("Parser id hash: ", hash)
	return hash
}

func (parser Parser) Parse(scanResult *src.ScanResult, config *cnf.Config) *src.ParseResult {

	logrus.Debug("Starting to parse")
	pr := src.ParseResult{
		FilePath: scanResult.FilePath,
	}

	id, err := machineid.ID()
	if err != nil {
		logrus.Error(generateIdError)
		pr.Errors = append(pr.Errors, generateIdError)
		pr.ParseId = parser.GenerateParseId(&pr, config)
		return &pr
	}
	pr.MachineId = id

	if len(scanResult.Books) == 0 {
		logrus.Error(emptyScanBookCountError)
		pr.Errors = append(pr.Errors, emptyScanBookCountError)
		pr.ParseId = parser.GenerateParseId(&pr, config)
		return &pr
	}

	var wg sync.WaitGroup

	booksCount := len(scanResult.Books)
	logrus.Debug("Books count: ", booksCount)
	wg.Add(booksCount)
	booksChan := make(chan src.Book, booksCount)
	errorsChan := make(chan error, booksCount)

	startTime := time.Now()
	pr.StartTime = startTime.Unix()

	for _, bookFile := range scanResult.Books {
		go parseBook(&wg, bookFile, config, booksChan, errorsChan)
	}

	logrus.Info("Waiting for parsing workers to finish")
	wg.Wait()
	close(booksChan)
	close(errorsChan)

	for el := range booksChan {
		pr.Books = append(pr.Books, el)
	}
	for el := range errorsChan {
		pr.Errors = append(pr.Errors, el)
	}

	elapsed := time.Since(startTime)
	pr.Duration = elapsed

	pr.ParseId = parser.GenerateParseId(&pr, config)
	logrus.Info("Parsing workers finished. Elapsed time: ", elapsed)
	logrus.Info("ParsedBooks: ", len(pr.Books), ". Errors: ", len(pr.Errors))

	return &pr
}

func HandleResult(handler src.ParseResultHandler, parserResult *src.ParseResult) bool {
	return handler.Handle(parserResult)
}

func parseBook(wg *sync.WaitGroup, bookFile src.BookFile, config *cnf.Config, bookChan chan src.Book, errorChan chan error) {
	logrus.Trace("Worker started for ", bookFile.FilePath)
	var (
		bookInfo *src.BookInfo
		err      error
	)

	defer func() {
		logrus.Trace("Worker finished for ", bookFile.FilePath)
		wg.Done()
	}()

	if bookFile.IsSymLink {
		errorChan <- src.ParseError{PreviousError: symlinkError, FileName: bookFile.FilePath}
		return
	}

	if !utils.IsStringInSlice(bookFile.Ext, config.ScanExt) {
		logrus.Debug("Skipped ext,", config.ScanExt, " for file ", bookFile.FilePath)
		return
	}

	switch bookFile.Ext {
	case extPDF:
		parser := pdf_parser.PdfParser{}
		bookInfo, err = parser.Parse(&bookFile, config.WithCoverImages)
	case extFB2:
		parser := fb2_parser.Fb2Parser{}
		bookInfo, err = parser.Parse(&bookFile, config.WithCoverImages)
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

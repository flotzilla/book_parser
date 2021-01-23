package src

import (
	cnf "book_parser/src/config"
	"time"
)

type BookFile struct {
	FilePath  string
	Name      string
	FileSize  int64
	Ext       string
	ModTime   time.Time
	IsSymLink bool
}

type BookInfo struct {
	Title         string
	Author        string
	Authors       []string
	Tags          string // coma separated values
	ISBN          string
	NumberOfPages int
	Description   string
	Genres        []string
	CoverPage     CoverPage
	Language      string
	Date          string
	PublisherInfo PublisherInfo
}

type CoverPage struct {
	ContentType string
	Data        string
}
type PublisherInfo struct {
	BookName  string
	Publisher string
	City      string
	Year      string
	ISBN      string
}

type Book struct {
	BookFile BookFile
	BookInfo BookInfo
}

type ConfigInterface interface {
	GetConfigHash() string
	ShowConfig()
}

type ScanResult struct {
	FilePath             string
	BooksFoundTotalCount int
	BooksSkippedCount    int
	BooksTotalCount      int

	Books []BookFile
}

type ParseResult struct {
	MachineId string        `json:"machineID"`
	ParseId   string        `json:"parseID"`
	FilePath  string        `json:"filePath"`
	StartTime int64         `json:"startTime"`
	Duration  time.Duration `json:"duration"`
	Books     []Book        `json:"book"`
	Errors    []error       `json:"error"`
}

// TODO rework handler
type ParseResultHandler interface {
	Handle(result *ParseResult) bool
}

type ParserInterface interface {
	Parse(scanResult *ScanResult, config *cnf.Config) *ParseResult
	GenerateParseId(result *ParseResult, config *cnf.Config) string
}

type ParseError struct {
	PreviousError error
	FileName      string
}

func (p ParseError) Error() string {
	panic(p.PreviousError.Error() + ". Filename: " + p.FileName)
}

type ParserInfoInterface interface {
	Parse(bookFile *BookFile, withCover bool) (*BookInfo, error)
}

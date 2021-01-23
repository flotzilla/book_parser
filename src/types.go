package src

import "time"

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

type BookParser interface {
	Parse(bookFile *BookFile) *BookInfo
}

type BookGrabber interface {
	GrabData(bookInfo *BookInfo)
}

type Config struct {
	ScanExt         []string
	SkippedExt      []string
	WithCoverImages bool
}

type ScanResult struct {
	BooksFoundTotalCount int
	BooksSkippedCount    int
	BooksTotalCount      int

	Books []BookFile
}

type ParseResult struct {
	Books     []Book
	Errors    []error
	StartTime int64
	Duration  time.Duration
}

type ParseResultHandler interface {
	Handle(result *ParseResult) bool
}

type ParseError struct {
	PreviousError error
	FileName      string
}

func (p ParseError) Error() string {
	panic(p.PreviousError.Error() + ". Filename: " + p.FileName)
}

type Parser interface {
	Parse(bookFile *BookFile, withCover bool) (*BookInfo, error)
}

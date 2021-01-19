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
	Pages         int
	ISBN          string
	NumberOfPages int
	Subject       string
	Description   string
	Genres        []string
	CoverPage     string
	Language      string
	Date          string
	Publisher     Publisher
	PublisherInfo PublisherInfo
}

type Publisher struct {
	FirstName string
	LastName  string
	id        string
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
	Books  []Book
	Errors []error
	// time start
	// time end
	// books counter ??
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

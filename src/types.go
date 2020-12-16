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
	Pages         int
	ISBN          string
	NumberOfPages int
	Subject       string
	Description   string
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
	ScanExt    []string
	SkippedExt []string
	HasInit    bool
}

type ScanResult struct {
	BooksFoundTotalCount int
	BooksSkippedCount    int
	BooksTotalCount      int

	Books []BookFile
}

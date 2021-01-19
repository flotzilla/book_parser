package fb2_parser

import (
	"book_parser/src"
	fb2_parser "book_parser/src/parser/types/fb2"
	"testing"
	"time"
)

func TestParseFb2WithoutCover(t *testing.T) {
	parser := fb2_parser.Fb2Parser{}

	b := src.BookFile{
		FilePath:  "book.fb2",
		Name:      "test",
		FileSize:  0,
		Ext:       "",
		ModTime:   time.Time{},
		IsSymLink: false,
	}
	result, err := parser.Parse(&b, false)

	if err == nil {
		t.Error(err)
	}

	if result.Title != "Test book" {
		t.Errorf("Invalid title parsing")
	}
}

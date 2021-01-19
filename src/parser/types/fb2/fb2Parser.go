package fb2_parser

import (
	"book_parser/src"
	"strings"
)

type Fb2Parser struct{}

func (parser *Fb2Parser) Parse(bookFile *src.BookFile, withCover bool) (*src.BookInfo, error) {
	fb2Book, err := parseFb2File(bookFile.FilePath, withCover)

	if err != nil {
		return nil, err
	}

	bI := src.BookInfo{
		Title:         fb2Book.GetTitle(),
		Authors:       fb2Book.GetAuthors(),
		Pages:         fb2Book.GetPages(),
		ISBN:          fb2Book.GetISBN(),
		Subject:       fb2Book.GetKeywords(),
		Description:   fb2Book.GetAnnotation(),
		Genres:        fb2Book.GetGenres(),
		Language:      fb2Book.GetLanguage(),
		Date:          fb2Book.GetDate(),
		PublisherInfo: fb2Book.PublisherInfo,
		Publisher:     fb2Book.Publisher,
	}

	if withCover {
		bI.CoverPage = fb2Book.Cover
	}

	return &bI, nil
}

func parseFb2File(filePath string, withCover bool) (*Fb2Info, error) {
	return &Fb2Info{}, nil
}

type Fb2Info struct {
	Title         string
	Authors       []string
	Pages         string
	ISBN          string
	Annotation    string
	Genres        []string
	Keywords      string
	Language      string
	Cover         string
	Date          string
	Publisher     src.Publisher
	PublisherInfo src.PublisherInfo
}

func (book *Fb2Info) GetTitle() string {
	return book.Title
}

func (book *Fb2Info) GetAuthors() []string {
	return book.Authors
}

func (book *Fb2Info) GetAuthorsAsString() string {
	return strings.Join(book.Authors, ", ")
}

func (book *Fb2Info) GetPages() int {
	return 0
}

func (book *Fb2Info) GetISBN() string {
	return book.ISBN
}

func (book *Fb2Info) GetAnnotation() string {
	return book.Annotation
}

func (book *Fb2Info) GetGenres() []string {
	return book.Genres
}

func (book *Fb2Info) GetKeywords() string {
	return book.Keywords
}

func (book *Fb2Info) GetLanguage() string {
	return book.Language
}

func (book *Fb2Info) GetCover() string {
	return book.Cover
}

func (book *Fb2Info) GetDate() string {
	return book.Date
}

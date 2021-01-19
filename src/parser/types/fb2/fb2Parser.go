package fb2_parser

import (
	"book_parser/src"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Fb2Parser struct{}

func (parser *Fb2Parser) Parse(bookFile *src.BookFile, withCover bool) (*src.BookInfo, error) {
	fb2Book, err := parseFb2File(bookFile.FilePath)

	if err != nil {
		return nil, err
	}

	bI := src.BookInfo{
		Title:         fb2Book.Description.TitleInfo.BookTitle,
		Authors:       humanDataToStringArray(&fb2Book.Description.TitleInfo.Author),
		ISBN:          fb2Book.Description.PublishInfo.ISBN,
		NumberOfPages: 0,
		Tags:          fb2Book.Description.TitleInfo.Keywords,
		Description:   fb2Book.Description.TitleInfo.Annotation,
		Genres:        fb2Book.Description.TitleInfo.Genre,
		Language:      fb2Book.Description.TitleInfo.Lang,
		Date:          fb2Book.Description.TitleInfo.Date,
		PublisherInfo: src.PublisherInfo{
			BookName:  fb2Book.Description.PublishInfo.BookName,
			Publisher: fb2Book.Description.PublishInfo.Publisher,
			City:      fb2Book.Description.PublishInfo.City,
			Year:      fb2Book.Description.PublishInfo.Year,
			ISBN:      fb2Book.Description.PublishInfo.ISBN,
		},
	}

	if withCover && fb2Book.Description.TitleInfo.CoverPage.Image.Href != "" {
		c, err := fb2Book.getCover()
		if err != nil {
			return nil, err
		}

		bI.CoverPage = *c
	}
	return &bI, nil
}

func parseFb2File(filePath string) (*Fb2Info, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	var book Fb2Info

	err = xml.Unmarshal(bytes, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

type Fb2Info struct {
	FictionBook xml.Name `xml:"FictionBook"`
	Description struct {
		TitleInfo struct {
			Genre      []string    `xml:"genre"`
			GenreType  []string    `xml:"genre-type"`
			Author     []HumanData `xml:"author"`
			BookTitle  string      `xml:"book-title"`
			Annotation string      `xml:"annotation"`
			Keywords   string      `xml:"keywords"`
			Date       string      `xml:"date"`
			CoverPage  struct {
				Image struct {
					Href string `xml:"href,attr"`
				} `xml:"image,allowempty"`
			} `xml:"coverpage"`
			Lang       string      `xml:"lang"`
			SrcLang    string      `xml:"src-lang"`
			Translator []HumanData `xml:"translator"`
		} `xml:"title-info"`
		DocumentInfo struct {
			Author      []HumanData `xml:"author"`
			ProgramUsed string      `xml:"program-used"`
			Date        string      `xml:"date"`
			Id          string      `xml:"id"`
			Version     string      `xml:"version"`
			History     string      `xml:"history"`
		} `xml:"document-info"`
		PublishInfo struct {
			BookName  string `xml:"book-name"`
			Publisher string `xml:"publisher"`
			City      string `xml:"city"`
			Year      string `xml:"year"`
			ISBN      string `xml:"isbn"`
		} `xml:"publish-info"`
	} `xml:"description"`
	Binary []struct {
		Id          string `xml:"id,attr"`
		ContentType string `xml:"content-type,attr"`
		Val         string `xml:",chardata"`
	} `xml:"binary"`
	// will skip body part, only parsing metadata and cover
}

type HumanData struct {
	FirstName  string `xml:"first-name"`
	MiddleName string `xml:"middle-name"`
	LastName   string `xml:"last-name"`
	NickName   string `xml:"nickname"`
}

func (b *Fb2Info) getCover() (*src.CoverPage, error) {
	str := b.Description.TitleInfo.CoverPage.Image.Href
	// remove # prefix from image href (# - for local content)
	id := str[1:len(str)]
	for _, el := range b.Binary {
		if el.Id == id {
			data, err := base64.StdEncoding.DecodeString(el.Val)
			if err != nil {
				return nil, err
			}

			cover := src.CoverPage{
				ContentType: el.ContentType,
				Data:        string(data),
			}

			return &cover, nil
		}
	}
	return nil, errors.New("cannot find document cover source")
}

func humanDataToStringArray(data *[]HumanData) []string {
	var humans []string
	for _, el := range *data {
		humans = append(humans, strings.Trim(
			el.FirstName+" "+el.MiddleName+" "+el.LastName+" "+el.NickName,
			" "))
	}

	return humans
}

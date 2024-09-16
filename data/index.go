package data

import "github.com/man-on-box/litepage/util"

type Head struct {
	Title string
}

type PageIndex struct {
	Head
	Header    string
	Subheader string
	DocsUrl   string
}

type HomepageData struct {
	Title     string `json:"title"`
	Header    string `json:"header"`
	Subheader string `json:"subheader"`
	DocsUrl   string `json:"docsUrl"`
}

func (d *Data) NewPageIndex() PageIndex {
	data := HomepageData{}
	util.ParseJSONFile("content/homepage.json", &data)

	return PageIndex{
		Head: Head{
			Title: data.Title,
		},
		Header:    data.Header,
		Subheader: data.Subheader,
		DocsUrl:   data.DocsUrl,
	}
}

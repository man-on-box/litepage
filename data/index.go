package data

type Head struct {
	Title string
}

type PageIndex struct {
	Head
}

func (d *Data) NewPageIndex() PageIndex {
	return PageIndex{
		Head: Head{
			Title: "Litepage app",
		},
	}
}

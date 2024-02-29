package dto

type TwonkyAlbumList struct {
	Channel struct {
		AlbumArtURI         string `xml:"albumArtURI"`
		ChildCount          int    `xml:"childCount"`
		ChildCountContainer int    `xml:"childCountContainer"`
		Class               string `xml:"class"`
		ContainerContent    string `xml:"containerContent"`
		Copyright           string `xml:"copyright"`
		Description         string `xml:"description"`
		ID                  string `xml:"id"`
		Item                []struct {
			Bookmark  string `xml:"bookmark"`
			Enclosure struct {
				Type string `xml:"type,attr"`
				Url  string `xml:"url,attr"`
			} `xml:"enclosure"`
			Meta struct {
				ChildCount  int    `xml:"childCount,attr"`
				ID          string `xml:"id,attr"`
				ParentID    string `xml:"parentID,attr"`
				Restricted  int    `xml:"restricted,attr"`
				Searchable  int    `xml:"searchable,attr"`
				Sec         string `xml:"sec,attr"`
				AlbumArtURI struct {
					ProfileID string `xml:"profileID,attr"`
					CharData  string `xml:",chardata"`
				} `xml:"albumArtURI"`
				ChildCountContainer int    `xml:"childCountContainer"`
				Class               string `xml:"class"`
				ContainerContent    string `xml:"containerContent"`
				ModificationTime    int    `xml:"modificationTime"`
				Title               string `xml:"title"`
			} `xml:"meta"`
			Title string `xml:"title"`
		} `xml:"item"`
		Language         string `xml:"language"`
		Link             string `xml:"link"`
		ModificationTime int    `xml:"modificationTime"`
		ParentList       struct {
			Parent []struct {
				ChildCount int    `xml:"childCount"`
				Class      string `xml:"class"`
				ID         string `xml:"id"`
				Title      string `xml:"title"`
				Url        string `xml:"url"`
			} `xml:"parent"`
		} `xml:"parentList"`
		PubDate       string `xml:"pubDate"`
		Returneditems string `xml:"returneditems"`
		Title         string `xml:"title"`
		Url           string `xml:"url"`
	} `xml:"channel"`
}

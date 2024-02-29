package dto

type TwonkySongList struct {
	Channel struct {
		AlbumArtURI         string               `xml:"albumArtURI"`
		ChildCount          int                  `xml:"childCount"`
		ChildCountContainer int                  `xml:"childCountContainer"`
		Class               string               `xml:"class"`
		ContainerContent    string               `xml:"containerContent"`
		Copyright           string               `xml:"copyright"`
		Description         string               `xml:"description"`
		ID                  string               `xml:"id"`
		Item                []TwonkySongListItem `xml:"item"`
		Language            string               `xml:"language"`
		Link                string               `xml:"link"`
		ModificationTime    int                  `xml:"modificationTime"`
		ParentList          struct {
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

type TwonkySongListItem struct {
	Bookmark  string `xml:"bookmark"`
	Enclosure struct {
		Type string `xml:"type,attr"`
		Url  string `xml:"url,attr"`
	} `xml:"enclosure"`
	Meta struct {
		ID          string `xml:"id,attr"`
		ParentID    string `xml:"parentID,attr"`
		RefID       string `xml:"refID,attr"`
		Restricted  int    `xml:"restricted,attr"`
		Sec         string `xml:"sec,attr"`
		AddedTime   int    `xml:"addedTime"`
		Album       string `xml:"album"`
		AlbumArtURI struct {
			ProfileID string `xml:"profileID,attr"`
			CharData  string `xml:",chardata"`
		} `xml:"albumArtURI"`
		AlbumArtist         string  `xml:"albumArtist"`
		Artist              string  `xml:"artist"`
		Class               string  `xml:"class"`
		Creator             string  `xml:"creator"`
		Date                string  `xml:"date"`
		Duration            string  `xml:"duration"`
		Extension           string  `xml:"extension"`
		Format              string  `xml:"format"`
		Genre               string  `xml:"genre"`
		LastPlayedTime      *string `xml:"lastPlayedTime"`
		ModificationTime    int     `xml:"modificationTime"`
		NumberOfThisDisc    int     `xml:"numberOfThisDisc"`
		OriginalTrackNumber int     `xml:"originalTrackNumber"`
		PersistentBookmark  string  `xml:"PersistentBookmark"`
		Playcount           *int    `xml:"playcount"`
		Res                 struct {
			Bitrate         int    `xml:"bitrate,attr"`
			Duration        string `xml:"duration,attr"`
			ProtocolInfo    string `xml:"protocolInfo,attr"`
			SampleFrequency int    `xml:"sampleFrequency,attr"`
			Size            int    `xml:"size,attr"`
			Timeseekinfo    string `xml:"timeseekinfo,attr"`
			CharData        string `xml:",chardata"`
		} `xml:"res"`
		Title string `xml:"title"`
	} `xml:"meta"`
	Title string `xml:"title"`
}

package dto

type TwonkyGengeList struct {
	ChildCount  string `json:"childCount"`
	Copyright   string `json:"copyright"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Item        []struct {
		Bookmark  string `json:"bookmark"`
		Enclosure struct {
			Type  string `json:"type"`
			URL   string `json:"url"`
			Value string `json:"value"`
		} `json:"enclosure"`
		Meta struct {
			ChildCount             string `json:"childCount"`
			Dc_Title               string `json:"dc:title"`
			ID                     string `json:"id"`
			ParentID               string `json:"parentID"`
			Pv_ChildCountContainer string `json:"pv:childCountContainer"`
			Pv_ContainerContent    string `json:"pv:containerContent"`
			Pv_ModificationTime    string `json:"pv:modificationTime"`
			Restricted             string `json:"restricted"`
			Searchable             string `json:"searchable"`
			Upnp_AlbumArtURI       string `json:"upnp:albumArtURI"`
			Upnp_Class             string `json:"upnp:class"`
			Xmlns_Sec              string `json:"xmlns:sec"`
		} `json:"meta"`
		Title string `json:"title"`
	} `json:"item"`
	Language   string `json:"language"`
	Link       string `json:"link"`
	ParentList []struct {
		ChildCount string `json:"childCount"`
		ID         string `json:"id"`
		Title      string `json:"title"`
		Upnp_Class string `json:"upnp:class"`
		URL        string `json:"url"`
	} `json:"parentList"`
	PubDate                string `json:"pubDate"`
	Pv_ChildCountContainer string `json:"pv:childCountContainer"`
	Pv_ContainerContent    string `json:"pv:containerContent"`
	Pv_ModificationTime    string `json:"pv:modificationTime"`
	Returneditems          string `json:"returneditems"`
	Title                  string `json:"title"`
	Upnp_Class             string `json:"upnp:class"`
	URL                    string `json:"url"`
}

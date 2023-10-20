package dto

type ImportSongs struct {
	Path          string   `json:"path"`
	Recursive     bool     `json:"recursive"`
	SongExtension []string `json:"songExtension"`
}

package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"
	"playlistturbo.com/config"
	"playlistturbo.com/plterror"
)

type Utils interface {
	Encoder
	Param
	Handler
}

type utils struct{}

func NewUtils() Utils {
	return utils{}
}

func PrintMetadata(m tag.Metadata) {
	fmt.Printf("Metadata Format: %v\n", m.Format())
	fmt.Printf("File Type: %v\n", m.FileType())

	fmt.Printf(" Title: %v\n", m.Title())
	fmt.Printf(" Album: %v\n", m.Album())
	fmt.Printf(" Artist: %v\n", m.Artist())
	fmt.Printf(" Composer: %v\n", m.Composer())
	fmt.Printf(" Genre: %v\n", m.Genre())
	fmt.Printf(" Year: %v\n", m.Year())

	track, trackCount := m.Track()
	fmt.Printf(" Track: %v of %v\n", track, trackCount)

	disc, discCount := m.Disc()
	fmt.Printf(" Disc: %v of %v\n", disc, discCount)

	fmt.Printf(" Picture: %v\n", m.Picture())
	fmt.Printf(" Lyrics: %v\n", m.Lyrics())
	fmt.Printf(" Comment: %v\n", m.Comment())
}

func GetMp3Time(r *os.File) float64 {
	t := 0.0
	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			fmt.Println(err, r.Name())
			return t
		}
		t = t + f.Duration().Seconds()
	}

	return t
}

func SecondsToMinutes(inSeconds int) string {
	out := ""
	if inSeconds > 3599 {
		hours := inSeconds / 3600
		minutes := inSeconds / 60
		seconds := inSeconds % 60
		out = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		minutes := inSeconds / 60
		seconds := inSeconds % 60
		out = fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
	return out
}

func ValidSongExtension(extension string) error {
	for _, ext := range config.Config.SupportedExtensions {
		if ext == extension {
			return nil
		}
	}
	return plterror.InvalidExtension
}

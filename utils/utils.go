package utils

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

func RemoveAccent(word string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, word)
	return result
}

func ExtractYear(s string) int {
	year := 0
	re := regexp.MustCompile("[0-9]+")
	arrNumber := re.FindAllString(s, -1)

	for _, v := range arrNumber {
		val, _ := strconv.Atoi(v)
		if val > 1900 && val < 3000 {
			year = val
		}
	}
	return year
}

func ExtractPath(path string) string {
	temp := strings.Split(path, "/")
	var newPath string
	for x := 1; x < len(temp)-1; x++ {
		piece := "/" + temp[x]
		newPath += piece
	}
	return newPath
}

func ExtractFilename(path string) string {
	pathSplit := strings.Split(path, "/")
	pIndex := len(pathSplit) - 1
	fname := pathSplit[pIndex]
	return fname
}

func ExtractAlbumName(s string) string {
	album := ""
	re := regexp.MustCompile(".*?[[(](\\d{4})[])]")
	ar1 := re.FindAllString(s, -1)
	if len(ar1) == 0 {
		re = regexp.MustCompile(".*?(\\d{4})")
		ar1 = re.FindAllString(s, -1)
	}
	if len(ar1) > 0 {
		ar2 := strings.TrimLeft(s, ar1[0])
		ar3 := strings.TrimLeft(ar2, " - ")
		ar4 := strings.TrimLeft(ar3, "-")
		album = strings.TrimSpace(ar4)
	} else {
		album = s
	}
	return album
}

func ExtractTitle(s string) string {
	title := ""
	re := regexp.MustCompile(".*?(\\d{2})")
	ar1 := re.FindAllString(s, -1)

	if len(ar1) > 0 {
		ar2 := strings.TrimLeft(s, ar1[0])
		ar3 := strings.TrimLeft(ar2, " - ")
		ar4 := strings.TrimLeft(ar3, "-")
		ar5 := strings.Split(ar4, ".mp3")
		title = strings.TrimSpace(ar5[0])
	} else {
		title = s
	}
	return title
}

func ValidateGAA(s string) bool {
	if s == "mp3tag" || s == "folder" {
		return true
	}
	return false
}

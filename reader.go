package oku

import (
	"fmt"
	"io"

	"github.com/qiniu/iconv"
	"github.com/saintfish/chardet"
)

// validEncodings is the intersection of types supported by chardet and iconv (ISO-8859-8-I is the only format not recognised by iconv)
var validEncodings = []string{
	"Big5",
	"EUC-JP", "EUC-KR",
	"ISO-2022-JP", "ISO-2022-KR", "ISO-2022-CN",
	"ISO-8859-1", "ISO-8859-2", "ISO-8859-5", "ISO-8859-6", "ISO-8859-7", "ISO-8859-8", "ISO-8859-9",
	"GB18030",
	"windows-1250", "windows-1251", "windows-1252", "windows-1253", "windows-1254", "windows-1255", "windows-1256",
	"KOI8-R",
	"Shift_JIS",
	"UTF-8", "UTF-16BE", "UTF-16LE", "UTF-32BE", "UTF-32LE",
}

type UTF8ReadCloser struct {
	r *iconv.Reader
	c iconv.Iconv
}

func (u *UTF8ReadCloser) Close() error {
	return u.c.Close()
}

func (u *UTF8ReadCloser) Read(b []byte) (int, error) {
	return u.r.Read(b)
}

func NewUTF8ReadCloser(r io.ReadSeeker, encoding string) (io.ReadCloser, error) {
	// rewind Seeker after encoding detection
	defer r.Seek(0, 0)

	// validate encoding for reading
	var found bool
	for _, v := range validEncodings {
		if v == encoding {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf(`detected file encoding:"%s", but there is no valid reader`, encoding)
	}

	cd, err := iconv.Open("utf-8", encoding)
	if err != nil {
		return nil, err
	}

	return &UTF8ReadCloser{r: iconv.NewReader(cd, r, 0), c: cd}, nil
}

type Detected struct {
	Charset    string
	Confidence int
}

func DetectEncoding(b []byte) (Detected, error) {
	d := chardet.NewTextDetector()
	res, err := d.DetectBest(b)
	if err != nil {
		return Detected{}, err
	}

	if res.Charset == "GB-18030" {
		// set canonical name for this encoding type (this is a chardet bug)
		res.Charset = "GB18030"
	}

	return Detected{Charset: res.Charset, Confidence: res.Confidence}, nil
}

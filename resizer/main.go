package resizer

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
	"path"
	"fmt"
	"net/http"
	"encoding/base64"
	"bytes"
	"strings"
	"io/ioutil"
)

type Resizer struct {
	fileRoot string
}

func NewResizer(root string) *Resizer {
	resizer := &Resizer{}

	resizer.fileRoot = root

	return resizer
}

func (r *Resizer) DownloadImage(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	s := strings.Split(url, "/")
	filename := s[len(s) - 1]
	
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	p := path.Join(r.fileRoot, filename)
	mode := int(0777)
	ioutil.WriteFile(p, body, os.FileMode(mode))

	return filename, nil
}

func (r *Resizer) ResizeImage(filePath string) (string, error) {
	// open "test.jpg"
	p := path.Join(r.fileRoot, filePath)
	file, err := os.Open(p)
	if err != nil {
		return "", err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return "", err
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(128, 128, img, resize.Lanczos3)

	np := path.Join(r.fileRoot, "test_resized.jpg")
	out, err := os.Create(np)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	// alright pass back the base 64 data
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(buf.Bytes())), nil
}
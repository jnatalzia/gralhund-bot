package resizer

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
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
	p := path.Join(r.fileRoot, filePath)
	file, err := os.Open(p)
	if err != nil {
		return "", err
	}

	imgType := strings.Split(filePath, ".")
	extension := imgType[len(imgType) - 1]

	var img image.Image
	// decode jpeg into image.Image
	if extension == "png" {
		img, err = png.Decode(file)	
	} else {
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		return "", err
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(128, 128, img, resize.Lanczos3)
	// alright pass back the base 64 data
	buf := new(bytes.Buffer)
	if extension == "png" {
		err = png.Encode(buf, m)
	} else {
		err = jpeg.Encode(buf, m, nil)
	}

	if err != nil {
		return "", err
	}

	return "data:image/" + extension + ";base64," + base64.StdEncoding.EncodeToString([]byte(buf.Bytes())), nil
}
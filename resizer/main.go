package resizer

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jnatalzia/gralhund-bot/utils"
	"github.com/nfnt/resize"
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
	// User-Agent:
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:66.0) Gecko/20100101 Firefox/66.0")
	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	s := strings.Split(url, "/")
	filename := s[len(s)-1]

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
	// TODO: Change file into bytearray early to prevent (presumed) double open
	mime, extension, err := mimetype.DetectFile(p)
	if err != nil {
		return "", err
	}
	fmt.Println("File info")
	fmt.Println(mime)

	allowedExtensions := []string{"png", "jpeg", "jpg"}

	if !utils.Contains(allowedExtensions, extension) {
		return "", errors.New("True file extension was " + extension + ". Must be one of png, jpeg, jpg")
	}

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

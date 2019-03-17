package resizer

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
	"path"
	"encoding/base64"
	"bytes"
)

type Resizer struct {
	fileRoot string
}

func NewResizer(root string) *Resizer {
	resizer := &Resizer{}

	resizer.fileRoot = root

	return resizer
}

func (r *Resizer) ResizeImage(filePath string) (string) {
	// open "test.jpg"
	p := path.Join(r.fileRoot, filePath)
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Thumbnail(128, 128, img, resize.Lanczos3)

	np := path.Join(r.fileRoot, "test_resized.jpg")
	out, err := os.Create(np)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	// alright pass back the base 64 data
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(buf.Bytes()))
}
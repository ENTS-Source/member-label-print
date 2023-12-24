package printer

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var SerialPort string

var img image.Image

func LoadLogo(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err = image.Decode(f)
	if err != nil {
		return err
	}

	return nil
}

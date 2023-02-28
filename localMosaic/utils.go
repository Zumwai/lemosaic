package localMosaic

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func leconvert(src image.Image) *image.RGBA {
	b := src.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
	return m
}

func getDecodedFile(name string) (image.Image, error) {
	/* file path checker needed? */
	src, err := ConvertCorrectType(name)

	if err != nil {
		return nil, err
	}
	return src, nil
}

func encodeToFile(path, name, suffix string, dst image.Image) error {
	newFile, err := os.Create(path + name + suffix)
	if err != nil {
		return err
	}
	defer newFile.Close()

	err = png.Encode(newFile, dst)
	if err != nil {
		return err
	}
	return nil
}

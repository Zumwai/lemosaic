package localMosaic

import (
	"image"
	"image/png"
	"os"
)

func getDecodedFile(name string) (*image.NRGBA, error) {
	/* file path checker needed? */
	src, err := ConvertCorrectType(name)

	if err != nil {
		return nil, err
	}
	return src, nil
}

func encodeToFile(path, name, suffix string, dst *image.NRGBA) error {
	newFile, err := os.Create(path + "/" + name + suffix)
	if err != nil {
		return err
	}
	defer newFile.Close()
	enc := png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
	err = enc.Encode(newFile, dst)
	if err != nil {
		return err
	}
	return nil
}

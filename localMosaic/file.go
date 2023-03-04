package localMosaic

import (
	//"image"
	"image/png"
	"mosaic/imgConv"
	"os"
)

/* decodes given image */
func getDecodedFile(name string) (imgConv.Image, error) {
	/* file path checker needed? */
	ret, err := ConvertCorrectType(name)

	if err != nil {
		return ret, err
	}
	return ret, nil
}

/* encodes to local file */
func encodeToFile(path, name, suffix string, dst imgConv.Image) error {
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

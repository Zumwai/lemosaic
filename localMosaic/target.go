package localMosaic

import (
	"math"
	"mosaic/config"
	"mosaic/imgConv"
	//"mosaic/logRuntime"
	"path"
)

func ExecutePouring(name string, chunk int, goCount int) error {
	src, err := getDecodedFile(name)
	if err != nil {
		return err
	}
	dst := imgConv.PreparePouring(src, chunk, goCount)
	err = encodeToFile("./target/", path.Base(name), "_squared.png", dst)
	if err != nil {
		return err
	}
	return nil
}

func ExecuteMosaic(name string, chunk int, goCount int) error {
	src, err := getDecodedFile(name)
	if err != nil {
		return err
	}

	source, err := PopulateHashDir("./pics/", chunk)
	//logRuntime.PrintMemory("after populating hash\n")

	if err != nil {
		return err
	}
	dst := imgConv.PrepareMosaic(src, chunk, goCount, source)
	//logRuntime.PrintMemory("after mosaic\n")

	err = encodeToFile("./target/", path.Base(name), "_mosaic.png", dst)
	if err != nil {
		return err
	}
	return nil
}

/* calcultes average colors of given file, resized it in memory if requested*/
func CalcAverageColours(name string, size int) (pic imgConv.ImgInfo, err error) {
	img, err := getDecodedFile(name)
	if err != nil {
		return
	}
	if size == 0 {
		pic.Square = img
	} else {
		pic.Square, err = imgConv.ResizeInMemory(img, size, size, config.InterpolLookup())
		if err != nil {
			return pic, err
		}
	}

	for x := 0; x < pic.Square.Bounds().Max.X; x++ {
		for y := 0; y < pic.Square.Bounds().Max.Y; y++ {
			col := imgConv.RgbaToPixel(pic.Square.At(x, y).RGBA())
			pic.Av.R += col.R
			pic.Av.G += col.G
			pic.Av.B += col.B
			pic.Av.A += col.A

		}
	}
	imgArea := float64(pic.Square.Bounds().Max.X * pic.Square.Bounds().Max.Y)
	pic.Av.R = math.Round(pic.Av.R / imgArea)
	pic.Av.G = math.Round(pic.Av.G / imgArea)
	pic.Av.B = math.Round(pic.Av.B / imgArea)
	pic.Av.A = math.Round(pic.Av.A / imgArea)
	return
}

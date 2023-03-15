package localMosaic

import (
	"fmt"
	"mosaic/imgConv"
	//"mosaic/logRuntime"
	"mosaic/config"
	"path"
)

/* main func to get poured image from local file, calculates and prints file into ./target/ */
func ExecutePouring(name string) error {
	src, err := getDecodedFile(name)
	if err != nil {
		return err
	}
	dst := imgConv.PreparePouring(src)
	err = EncodeToFile("./target/", path.Base(name), "_squared", dst)
	if err != nil {
		return err
	}
	return nil
}

/* main func to get mosaic image from local file, calculates and prints file into ./target/ */
func ExecuteMosaic(name string) error {
	src, err := getDecodedFile(name)
	if err != nil {
		return err
	}

	source, err := PopulateHashDir(config.SrcImagesLookup())
	if err != nil {
		return err
	}
	if len(source) < 1 {
		return fmt.Errorf("no images in source dir")
	}

	dst := imgConv.PrepareMosaic(src, source)
	err = EncodeToFile("./target/", path.Base(name), "_mosaic", dst)
	if err != nil {
		return err
	}
	return nil
}

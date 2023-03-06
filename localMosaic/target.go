package localMosaic

import (
	"mosaic/imgConv"
	//"mosaic/logRuntime"
	"path"
)

/* main func to get poured image from local file, calculates and prints file into ./target/ */

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

/* main func to get mosaic image from local file, calculates and prints file into ./target/ */
func ExecuteMosaic(name string, chunk int, goCount int) error {
	src, err := getDecodedFile(name)
	if err != nil {
		return err
	}

	source, err := PopulateHashDir("./pics")
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

package localMosaic

import (
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

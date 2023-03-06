package localMosaic

import (
	"mosaic/config"
	"mosaic/imgConv"
	//"mosaic/logRuntime"
	"path"
)

/* main func to get poured image from local file, calculates and prints file into ./target/ */

func ExecutePouring(name string) error {
	src, err := getDecodedFile(name)
	if err != nil {
		return err
	}
	dst := imgConv.PreparePouring(src)
	err = encodeToFile("./target/", path.Base(name), "_squared", dst)
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

	source, err := PopulateHashDir("./pics")
	//logRuntime.PrintMemory("after populating hash\n")
	size := config.ChunkLookup()
	goCount := config.RoutineLookup()
	if err != nil {
		return err
	}
	dst := imgConv.PrepareMosaic(src, size, goCount, source)
	//logRuntime.PrintMemory("after mosaic\n")

	err = encodeToFile("./target/", path.Base(name), "_mosaic", dst)
	if err != nil {
		return err
	}
	return nil
}

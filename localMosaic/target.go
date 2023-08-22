package localMosaic

import (
	"fmt"
	"mosaic/imgConv"
	//"mosaic/logRuntime"
	"mosaic/config"
	"os"
	"path"
	"sync"
)

func PouringOnlyOne() error {

	return nil
}

/* main func to get poured image from local file, calculates and prints file into ./target/ */
func ExecutePouring(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	defer file.Close()
	if err != nil {
		return err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		var wg sync.WaitGroup
		dirname := name
		dirReader, err := os.ReadDir(dirname)
		if err != nil {
			return err
		}
		wg.Add(len(dirReader))
		for _, f := range dirReader {
			go func(pict string) error {
				defer wg.Done()
				src, err := getDecodedFile(pict)
				if err != nil {
					fmt.Println(err)
					return err
				}
				dst := imgConv.PreparePouring(src)
				err = EncodeToFile("./target/", path.Base(pict), "_squared", dst)
				if err != nil {
					fmt.Println(err)
					return err
				}
				return nil
			}(dirname + "/" + f.Name())
		}
		wg.Wait()
	} else {
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

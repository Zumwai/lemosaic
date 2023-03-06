package main

import (
	"flag"
	"fmt"
	//"log"
	"mosaic/config"
	"mosaic/localMosaic"
	"mosaic/logRuntime"
	//"net/http"
	//"github.com/pkg/profile"
	//_ "net/http/pprof"
	//"github.com/google/pprof/profile"
	//"mosaic/serve"
	"mosaic/serve"
)

func main() {
	flag.Int("interpol", 0, "Choose one out of the for methods for interpolation")
	flag.Int("format", 0, "choose internal format [0-1] - RGBA or NRGBA")
	flag.Int("encoder", 0, "choose final image format [0-3] - jpeg, png, tiff, gif")
	flag.Int("chunk", 20, "size of a square in image")
	flag.Int("routine", 10, "number of gourutine in use")

	/*sets additional tasks for changing size and format of the images*/
	var redoSmaller = flag.String("resize", "", "add this and foldder which you want to redo smallification")
	var pourTarget = flag.String("pour", "", "Pour this file into squares")
	var calcAverage = flag.String("average", "", "what are average colors of a picture")
	var populateHash = flag.String("populate", "", "populate db with hashes of images")
	var mosaic = flag.String("mosaic", "", "mosaic dat image")
	var browser = flag.Bool("serve", false, "enable -serve flag in order ot start a server")
	flag.Parse()

	//p := profile.Start(profile.CPUProfile)
	//	defer profile.Start(profile.CPUProfile).Stop()
	//configTask(*chunkSize, *methodFlag)

	config.RegDecoders()
	if *pourTarget != "" {
		err := localMosaic.ExecutePouring(*pourTarget)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *redoSmaller != "" {
		err := localMosaic.ChangeSrcsSize(*redoSmaller)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *calcAverage != "" {
		av, err := localMosaic.CalcAverageSrcColours(*calcAverage)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(av)
		}
	}
	if *populateHash != "" {
		hash, err := localMosaic.PopulateHashDir(*populateHash)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(hash)
		}
	}
	if *mosaic != "" {
		err := localMosaic.ExecuteMosaic(*mosaic)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *browser {
		serve.StartServer()
	}

	logRuntime.PrintMemory("at the end\n")
	//
	// p.Stop()
}

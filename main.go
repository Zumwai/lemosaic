package main

import (
	"flag"
	"fmt"
	//"log"
	"mosaic/config"
	"mosaic/localMosaic"
	"mosaic/logRuntime"
	//"net/http"
	"github.com/pkg/profile"
	//_ "net/http/pprof"
	//"github.com/google/pprof/profile"
	//"mosaic/serve"
	"mosaic/serve"
	"os"
)

func printHelp() {
	fmt.Printf(`available commands are:
	-interpol [0-3] - choose one out of 4 interpoltion types
	-chunk [1-99] - choose one delimeter for size
	-resize - resizes all images from [name] and puts them into ./smaller
	-target - converts image [name] to png and puts it into ./target
	-pour - square a circle
	-populate -  populate hash from dirs of png
	-mosaic - print target as a mosaic
	folder used for source image file is ./pics/, not customizable
	`)
	os.Exit(1)
}

func main() {
	flag.Int("interpol", 0, "Choose one out of the for methods for interpolation")
	var chunkSize = flag.Int("chunk", 20, "size of a square in image")
	/*sets additional tasks for changing size and format of the images*/
	var redoSmaller = flag.String("resize", "", "add this and foldder which you want to redo smallification")
	var pourTarget = flag.String("pour", "", "Pour this file into squares")
	var calcAverage = flag.String("average", "", "what are average colors of a picture")
	var populateHash = flag.String("populate", "", "populate db with hashes of images")
	var mosaic = flag.String("mosaic", "", "mosaic dat image")
	var routine = flag.Int("routine", 10, "number of gourutine in use")
	/*helper*/
	var help = flag.Bool("help", false, "THERE IS NO HELP")
	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}
	defer profile.Start().Stop()
	//defer profile.Start(profile.CPUProfile).Stop()
	//configTask(*chunkSize, *methodFlag)
	chunk := *chunkSize
	config.RegDecoders()
	if *pourTarget != "" {
		err := localMosaic.ExecutePouring(*pourTarget, chunk, *routine)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *redoSmaller != "" {
		err := localMosaic.ChangeSrcsSize(*redoSmaller, chunk)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *calcAverage != "" {
		av, err := localMosaic.CalcAverageColours(*calcAverage, chunk)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(av)
		}
	}
	if *populateHash != "" {
		hash, err := localMosaic.PopulateHashDir(*populateHash, chunk)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(hash)
		}
	}
	if *mosaic != "" {
		err := localMosaic.ExecuteMosaic(*mosaic, chunk, *routine)
		if err != nil {
			fmt.Println(err)
		}
	}
	serve.StartServer()
	logRuntime.PrintMemory("at the end\n")
}

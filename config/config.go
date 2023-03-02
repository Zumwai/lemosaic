package config

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

func staticInterpolLookupTable(method int) string {
	var methods = []string{
		"NearestNeighbor",
		"ApproxBiLinear",
		"CatmullRom",
		"BiLinear",
	}
	if method < 0 || method > 3 {
		fmt.Printf("there is only 4 types of interpolation available, choose %s, currently using default %s\n", methods, methods[0])
		return methods[0]
	}
	return methods[method]
}

func InterpolLookup() string {
	return staticInterpolLookupTable(flag.Lookup("interpol").Value.(flag.Getter).Get().(int))
}

func ChunkLookup() int {
	return flag.Lookup("chunk").Value.(flag.Getter).Get().(int)
}

func RoutineLookup() int {
	return flag.Lookup("routine").Value.(flag.Getter).Get().(int)
}

func RegDecoders() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

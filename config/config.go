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
	if method < 0 || method > len(methods)-1 {
		fmt.Printf("there is only %d types of interpolation available, choose %s, currently using default %s\n",
			len(methods), methods, methods[0])
		return methods[0]
	}
	return methods[method]
}

func InterpolLookup() string {
	return staticInterpolLookupTable(flag.Lookup("interpol").Value.(flag.Getter).Get().(int))
}

func staticFormatLookupTable(format int) string {
	var formats = []string{
		"RGBA",
		"NRGBA",
	}
	if format < 0 || format > len(formats)-1 {
		fmt.Printf("there is only %d types of interpolation available, choose %s, currently using default %s\n",
			len(formats), formats, formats[0])
		return formats[0]
	}
	return formats[format]
}

func FormatLookup() string {
	return staticFormatLookupTable(flag.Lookup("format").Value.(flag.Getter).Get().(int))
}

func ChunkLookup() int {
	return flag.Lookup("chunk").Value.(flag.Getter).Get().(int)
}

func staticEncoderLookup(encoder int) string {
	var encoders = []string{
		"jpeg",
		"png",
		"tiff",
		"gif",
	}
	if encoder < 0 || encoder > len(encoders)-1 {
		fmt.Printf("there is only %d types of interpolation available, choose %s, currently using default %s\n",
			len(encoders), encoders, encoders[0])
		return encoders[0]
	}
	return encoders[encoder]
}

func EncoderLookup() string {
	return staticEncoderLookup(flag.Lookup("encoder").Value.(flag.Getter).Get().(int))
}
func RoutineLookup() int {
	return flag.Lookup("routine").Value.(flag.Getter).Get().(int)
}

func RegDecoders() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

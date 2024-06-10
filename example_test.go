package osmshortlink

import (
	"fmt"
)

func ExampleCreate() {
	shortLink, err := Create(46.05141, 14.50604, 17)
	if err != nil {
		panic(err)
	}
	fmt.Println(shortLink)
	// Output: https://osm.org/go/0Ik3VNr_A-?m
}

func ExampleEncode() {
	shortLink, err := Encode(46.05141, 14.50604, 17)
	if err != nil {
		panic(err)
	}
	fmt.Println(shortLink)
	// Output: 0Ik3VNr_A-
}

func ExampleDecode() {
	latitude, longitude, zoom, err := Decode("0Ik3VNr_A-")
	if err != nil {
		panic(err)
	}
	fmt.Println(latitude, longitude, zoom)
	// Output: 46.05140447616577 14.506051540374756 17
}

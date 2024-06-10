package osmshortlink_test

import (
	"fmt"

	"github.com/stefanb/osmshortlink-go"
)

func ExampleCreate() {
	shortLink, err := osmshortlink.Create(46.05141, 14.50604, 17)
	if err != nil {
		panic(err)
	}
	fmt.Println(shortLink)
	// Output: https://osm.org/go/0Ik3VNr_A-?m
}

func ExampleEncode() {
	shortLink, err := osmshortlink.Encode(46.05141, 14.50604, 17)
	if err != nil {
		panic(err)
	}
	fmt.Println(shortLink)
	// Output: 0Ik3VNr_A-
}

func ExampleDecode() {
	latitude, longitude, zoom, err := osmshortlink.Decode("0Ik3VNr_A-")
	if err != nil {
		panic(err)
	}
	fmt.Println(latitude, longitude, zoom)
	// Output: 46.05140447616577 14.506051540374756 17
}

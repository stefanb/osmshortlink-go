// Creating OpenStreetMap short links, see https://wiki.openstreetmap.org/wiki/Shortlink
package osmshortlink

import (
	"fmt"
	"math"
	"strings"
)

const (
	BASE_SHORT_OSM_URL = "https://osm.org/go/"
)

// intToBase64 is a lookup table that translates 6-bit positive integer
// index values into their "Base64 Alphabet" equivalents as specified
// in Table 1 of RFC 2045.
var intToBase64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_~"

// Create a short link to OSM map with marker at he given position and zoom
func Create(latitude float32, longitude float32, zoom int) (string, error) {
	s, err := Encode(latitude, longitude, zoom)
	if err != nil {
		return "", err
	}
	return BASE_SHORT_OSM_URL + s + "?m", nil
}

// Encode a location and zoom, return a short string representing it.
func Encode(latitude float32, longitude float32, zoom int) (string, error) {
	if zoom < 0 || zoom > 20 {
		return "", fmt.Errorf("invalid zoom %d", zoom)
	}
	if latitude < -90 || latitude >= 90 {
		return "", fmt.Errorf("invalid latitude %v", latitude)
	}
	if longitude < -180 || longitude >= 180 {
		return "", fmt.Errorf("invalid longitude %v", longitude)
	}
	lat := uint32((latitude + 90) / 180 * (1 << 32))
	lon := uint32((longitude + 180) / 360 * (1 << 32))

	code := interleaveBits(lon, lat)
	str := ""
	for i := 0; i < int(math.Ceil((float64(zoom+8))/3)); i++ {
		str += string(intToBase64[int((code >> (58 - 6*i) & 0x3f))])
	}
	for j := 0; j < (zoom+8)%3; j++ {
		str += "-"
	}

	return str, nil
}

// interleaveBits combines two 32-bit integers into a 64-bit integer
func interleaveBits(x uint32, y uint32) uint64 {
	var c uint64 = 0
	for b := 31; b >= 0; b-- {
		c = c<<1 | uint64(x)>>b&1
		c = c<<1 | uint64(y)>>b&1
	}
	return c
}

// Decode a short string into a location and zoom.
func Decode(s string) (float64, float64, int, error) {
	if len(s) < 1 {
		return 0, 0, 0, fmt.Errorf("invalid osm short link string %q", s)
	}
	// convert old shortlink format to current one
	s = strings.ReplaceAll(s, "@", "~")
	var i int
	var x int64 = 0
	var y int64 = 0
	var z int = -8
	for i = 0; i < len(s); i++ {
		var digit int = -1
		var c byte = s[i]
		for j := 0; j < len(intToBase64); j++ {
			if c == intToBase64[j] {
				digit = j
				break
			}
		}
		if digit < 0 {
			break
		}
		if digit < 0 {
			break
		}
		// distribute 6 bits into x and y
		x <<= 3
		y <<= 3
		for j := 2; j >= 0; j-- {
			x |= ((map[bool]int64{true: 0, false: (1 << j)})[(digit&(1<<(j+j+1))) == 0])
			y |= ((map[bool]int64{true: 0, false: (1 << j)})[(digit&(1<<(j+j))) == 0])

		}
		z += 3
	}

	var lon float64 = float64(x)*math.Pow(float64(2), float64(2-3*i))*90.0 - float64(180)
	var lat float64 = float64(y)*math.Pow(float64(2), float64(2-3*i))*45.0 - float64(90)
	// adjust z
	if i < len(s) && s[i] == '-' {
		z -= 2
		if i+1 < len(s) && s[i+1] == '-' {
			z++
		}
	}
	// return new GeoParsedPoint(lat, lon, z);
	return lat, lon, z, nil
}

package osmshortlink

import (
	"math"
	"testing"
)

func Test_CreateOSMShortLinkFull(t *testing.T) {
	type args struct {
		latitude  float32
		longitude float32
		zoom      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"null island", args{0, 0, 0}, "https://osm.org/go/wAA--?m"},                   // http://www.salesianer.de/util/osmshortlinks.php
		{"null island", args{0, 0, 1}, "https://osm.org/go/wAA?m"},                     // http://www.salesianer.de/util/osmshortlinks.php
		{"null island", args{0, 0, 19}, "https://osm.org/go/wAAAAAAAA?m"},              // http://www.salesianer.de/util/osmshortlinks.php
		{"OsmAnd 3", args{40.59, -115.213, 9}, "https://osm.org/go/TelHTB--?m"},        // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L397
		{"Elixir 1", args{51.5110, 0.0550, 9}, "https://osm.org/go/0EEQjE--?m"},        // https://hex.pm/packages/osm_shortlink
		{"Slovenia rough ints", args{46, 14, 9}, "https://osm.org/go/0IO_cI--?m"},      // http://www.salesianer.de/util/osmshortlinks.php
		{"Slovenia rough float", args{46.0, 14.0, 9}, "https://osm.org/go/0IO_cI--?m"}, // http://www.salesianer.de/util/osmshortlinks.php
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Create(tt.args.latitude, tt.args.longitude, tt.args.zoom); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
				if err != nil {
					t.Errorf("Error: %v", err)
				}
			}
		})
	}
}

func Test_interleaveBits(t *testing.T) {
	type args struct {
		x uint32
		y uint32
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"zero", args{0, 0}, 0},
		{"one", args{1, 1}, 3},
		{"oneBinary", args{0b1, 0b1}, 0b0000000000000000000000000000000000000000000000000000000000000011},
		{"lastFirst", args{0b00000000000000000000000000000001, 0b10000000000000000000000000000000}, 0b100000000000000000000000000000000000000000000000000000000000010},
		{"firstLast", args{0b10000000000000000000000000000000, 0b00000000000000000000000000000001}, 0b1000000000000000000000000000000000000000000000000000000000000001},
		{"left", args{0b11111111111111111111111111111111, 0}, 0b1010101010101010101010101010101010101010101010101010101010101010},
		{"right", args{0, 0b11111111111111111111111111111111}, 0b0101010101010101010101010101010101010101010101010101010101010101},
		{"both", args{0b11111111111111111111111111111111, 0b11111111111111111111111111111111}, 0b1111111111111111111111111111111111111111111111111111111111111111},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := interleaveBits(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("interleaveBits() = %v (%b), want %v (%b)", got, got, tt.want, tt.want)
			}
		})
	}
}

func Fuzz_CreateDecodeShortLink(f *testing.F) {
	f.Add(float32(0), float32(0), 0)
	f.Add(float32(0), float32(0), 1)
	f.Add(float32(-90), float32(0), 15)
	f.Add(float32(90), float32(0), 15)
	f.Add(float32(0), float32(-180), 15)
	f.Add(float32(0), float32(180), 15)
	f.Add(float32(46.05141922831535), float32(14.506048858165741), 19)

	f.Fuzz(func(t *testing.T, origLat, origLon float32, origZoom int) {
		enc, encErr := Encode(origLat, origLon, origZoom)
		if origLat >= 90 || origLat < -90 || origLon >= 180 || origLon < -180 || origZoom > 20 || origZoom < 0 {
			if encErr == nil {
				t.Error("missing an encode error")
			}
			if enc != "" {
				t.Errorf("expected empty result, got %q", enc)
			}
			// return
		}

		decLat, decLon, decZoom, decErr := Decode(enc)
		if encErr != nil && decErr == nil {
			t.Errorf("missing a deocode error after %v", encErr)
		}

		if encErr == nil && decErr != nil {
			t.Errorf("decode error: %v", decErr)
			if decZoom != 0 {
				t.Errorf("Zoom decode with error shoould be 0, got after: %d, encoded: %q", decZoom, enc)
			}
			if decZoom != 0 {
				t.Errorf("Zoom decode with error shoould be 0, got after: %d, encoded: %q", decZoom, enc)
			}
			if decZoom != 0 {
				t.Errorf("Zoom decode with error shoould be 0, got after: %d, encoded: %q", decZoom, enc)
			}
		}

		if encErr == nil && decErr == nil {
			if origZoom != decZoom {
				t.Errorf("Zoom before: %d, after: %d, encoded: %q", origZoom, decZoom, enc)
			}

			tolerance := 0.8 // TODO: adjust it by zoom
			if math.Abs(float64(origLat)-decLat) > tolerance {
				t.Errorf("Lat before: %f, after: %f, delta %f, zoom %d, tolerance %f, encoded: %q", origLat, decLat, math.Abs(float64(origLat)-decLat), decZoom, tolerance, enc)
			}
			if math.Abs(float64(origLon)-decLon) > tolerance {
				t.Errorf("Lon before: %f, after: %f, delta %f, zoom %d, tolerance %f, encoded: %q", origLon, decLon, math.Abs(float64(origLon)-decLon), decZoom, tolerance, enc)
			}
		}

	})

}

func Test_CreateOSMShortLinkString(t *testing.T) {
	type args struct {
		latitude  float32
		longitude float32
		zoom      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"null island", args{0, 0, 0}, "wAA--"},      // http://www.salesianer.de/util/osmshortlinks.php
		{"null island", args{0, 0, 1}, "wAA"},        // http://www.salesianer.de/util/osmshortlinks.php
		{"null island", args{0, 0, 19}, "wAAAAAAAA"}, // http://www.salesianer.de/util/osmshortlinks.php
		// {"OsmAnd 1", args{51.51829, 0.07347, 16}, "0EEQsyfu"},   // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L395
		// {"OsmAnd 2", args{52.30103, 4.862927, 18}, "0E4_JiVhs"}, // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L396
		{"OsmAnd 3", args{40.59, -115.213, 9}, "TelHTB--"},        // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L397
		{"Elixir 1", args{51.5110, 0.0550, 9}, "0EEQjE--"},        // https://hex.pm/packages/osm_shortlink
		{"Slovenia rough ints", args{46, 14, 9}, "0IO_cI--"},      // http://www.salesianer.de/util/osmshortlinks.php
		{"Slovenia rough float", args{46.0, 14.0, 9}, "0IO_cI--"}, // http://www.salesianer.de/util/osmshortlinks.php
		// {"Trg republike", args{46.05, 14.5, 14}, "http://osm.org/go/0Ik1~phl-"},                              // http://www.salesianer.de/util/osmshortlinks.php
		// {"Prešernov trg", args{46.05141922831535, 14.506048858165741, 19}, "0Ik3VNupR"}, // http://www.salesianer.de/util/osmshortlinks.php
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Encode(tt.args.latitude, tt.args.longitude, tt.args.zoom); got != tt.want {
				t.Errorf("CreateOSMShortLinkString() = %v, want %v", got, tt.want)
				if err != nil {
					t.Errorf("Error: %v", err)
				}
			}
		})
	}
}

func TestDecodeShortLinkString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 float64
		want2 int
	}{
		{"null island", args{"wAA--"}, 0, 0, 0},      // http://www.salesianer.de/util/osmshortlinks.php
		{"null island", args{"wAA"}, 0, 0, 1},        // http://www.salesianer.de/util/osmshortlinks.php
		{"null island", args{"wAAAAAAAA"}, 0, 0, 19}, // http://www.salesianer.de/util/osmshortlinks.php
		// {"OsmAnd 1", args{51.51829, 0.07347, 16}, "0EEQsyfu"},   // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L395
		// {"OsmAnd 2", args{52.30103, 4.862927, 18}, "0E4_JiVhs"}, // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L396
		// {"OsmAnd 3", args{"TelHTB--"}, 40.59, -115.213, 9},        // https://github.com/osmandapp/OsmAnd/blob/485ac3690ce832fdd8ff96d303861162cb101570/OsmAnd-java/src/main/java/net/osmand/util/MapUtils.java#L397
		// {"Elixir 1", args{"0EEQjE--"}, 51.5110, 0.0550, 9},        // https://hex.pm/packages/osm_shortlink
		// {"Slovenia rough ints", args{"0IO_cI--"}, 46, 14, 9},      // http://www.salesianer.de/util/osmshortlinks.php
		// {"Slovenia rough float", args{"0IO_cI--"}, 46.0, 14.0, 9}, // http://www.salesianer.de/util/osmshortlinks.php
		// {"Trg republike", args{46.05, 14.5, 14}, "http://osm.org/go/0Ik1~phl-"},                              // http://www.salesianer.de/util/osmshortlinks.php
		// {"Prešernov trg", args{46.05141922831535, 14.506048858165741, 19}, "0Ik3VNupR"}, // http://www.salesianer.de/util/osmshortlinks.php

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := Decode(tt.args.s)
			if err != nil {
				t.Errorf("DecodeShortLinkString() got error %v", err)
			}
			if got != tt.want {
				t.Errorf("DecodeShortLinkString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeShortLinkString() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("DecodeShortLinkString() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

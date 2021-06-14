package colour

import (
	_ "embed"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"math"
)

//go:embed mappings.nbt
var mappingData []byte

// Block contains a name and properties value which can be used to find a block.
type Block struct {
	// Name is the name of the block.
	Name string `nbt:"name"`
	// Properties contains extra properties for the block, such as it's colour.
	Properties map[string]interface{} `nbt:"props"`
}

// Colour contains both an RGB value and a corresponding block.
type Colour struct {
	// Block is the block that matches the colour.
	Block Block `nbt:"block"`
	// RGB is the RGB of the colour.
	RGB `nbt:"rgb"`
}

// Defaults ...
func Defaults() (colours []Colour, err error) {
	err = nbt.Unmarshal(mappingData, &colours)
	if err != nil {
		return nil, err
	}
	return
}

// Closest finds the closest colour out of a list of colours.
func Closest(r, g, b uint32, colours []Colour) (match Colour) {
	minMSE := uint32(math.MaxUint32)

	for _, c := range colours {
		mse := uint32(c.ComputeMSE(int32(r), int32(g), int32(b)))
		if mse < minMSE {
			minMSE = mse
			match = c
		}
	}

	return match
}

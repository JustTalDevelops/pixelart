package commands

import (
	// These imports are to allow support for PNG, JPEG, and WEBP decoding.
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	_ "image/png"

	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/justtaldevelops/pixelart/colour"
	"github.com/nfnt/resize"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"image"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
)

// Image is a command which generates an image from a file in the working directory.
type Image struct {
	// Name is the name of the image.
	Name imageName `name:"name"`
	// ImageCreationDistance is how far away the image should be created from the player. This is 50 by default.
	ImageCreationDistance int `name:"distance" optional:""`
}

// Run is called when the command is ran from a source.
func (i Image) Run(source cmd.Source, output *cmd.Output) {
	if c, ok := source.(session.Controllable); ok {
		colours, err := colour.Defaults()
		if err != nil {
			output.Error(err)
			return
		}

		img := string(i.Name)

		reader, err := os.Open(img)
		if err != nil {
			output.Error(err)
			return
		}
		defer reader.Close()

		m, _, err := image.Decode(reader)
		if err != nil {
			output.Error(err)
			return
		}

		width, height := m.Bounds().Max.X, m.Bounds().Max.Y
		if height > cube.MaxY {
			scaleFactor := int(math.Ceil(float64(height / cube.MaxY))) + 1

			newWidth, newHeight := width/scaleFactor, height / scaleFactor
			m = resize.Resize(uint(newWidth), uint(newHeight), m, resize.Bicubic)

			width, height = newWidth, newHeight
		}

		facing := entity.Facing(c)
		spawnPos := cube.PosFromVec3(c.Position())

		if i.ImageCreationDistance == 0 {
			i.ImageCreationDistance = 50
		}

		for j := 0; j < i.ImageCreationDistance; j++ {
			spawnPos = spawnPos.Side(facing.Face())
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := m.At(x, y).RGBA()
				if a > 0 {
					var spawnX, spawnZ int

					switch facing {
					case cube.North:
						spawnX = spawnPos.X() + x - width / 2
						spawnZ = spawnPos.Z()
					case cube.South:
						spawnX = spawnPos.X() - x + width / 2
						spawnZ = spawnPos.Z()
					case cube.East:
						spawnX = spawnPos.X()
						spawnZ = spawnPos.Z() + x - width / 2
					case cube.West:
						spawnX = spawnPos.X()
						spawnZ = spawnPos.Z() - x + width / 2
					}

					closest := colour.Closest(r / 257, g / 257, b / 257, colours)
					pixel, _ := world.BlockByName(closest.Block.Name, closest.Block.Properties)

					c.World().SetBlock(cube.Pos{spawnX, height - y, spawnZ}, pixel)
				}
			}
		}

		output.Print(text.Colourf("<green>Successfully generated %v with blocks!</green>", img))
	} else {
		output.Error("This command can only be ran by a controllable!")
	}
}

// imageName is an enum which contains all images in the current working directory.
type imageName string

// Type is the type displayed client side.
func (imageName) Type() string {
	return "ImageName"
}

// Options contains each file that may be used to generate an image.
func (imageName) Options(_ cmd.Source) (opts []string) {
	files, _ := ioutil.ReadDir(".")

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" {
			opts = append(opts, file.Name())
		}
	}

	return
}

// SetOption ...
func (imageName) SetOption(option string, r reflect.Value) {
	r.SetString(option)
}

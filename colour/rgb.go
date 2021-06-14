package colour

// RGB is an additive color model in which red, green, and blue light are added together in various ways to reproduce a broad array of colors.
type RGB struct {
	// Red is the amount of red in the value.
	Red int32 `nbt:"red"`
	// Green is the amount of green in the value.
	Green int32 `nbt:"green"`
	// Blue is the amount of blue in the value.
	Blue int32 `nbt:"blue"`
}

// ComputeMSE ...
func (r RGB) ComputeMSE(pixR, pixG, pixB int32) int32 {
	return ((pixR-r.Red)*(pixR-r.Red) + (pixG-r.Green)*(pixG-r.Green) + (pixB-r.Blue)*(pixB-r.Blue)) / 3
}
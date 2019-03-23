package main

import (
	"fmt"
	"image/color"
	"math"
)

// Vec3 represents a 3D vector
type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func (v *Vec3) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}

// RGBA interpretation of the vector
func (v *Vec3) RGBA() color.Color {
	return color.RGBA{uint8(v.X * 255), uint8(v.Y * 255), uint8(v.Z * 255), 255}
}

// Length of the vector
func (v *Vec3) Length() float32 {
	return float32(math.Sqrt(float64(Dot(*v, *v))))
}

// Dot product between two vectors
func Dot(a Vec3, b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Normalize a vector
func Normalize(a Vec3) Vec3 {
	return MulScalar(1/a.Length(), a)
}

// Sub computes a - b
func Sub(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Add computes a + b
func Add(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// AddScalar computes s + a
func AddScalar(s float32, a Vec3) Vec3 {
	return Vec3{a.X + s, a.Y + s, a.Z + s}
}

// MulScalar computes s * a
func MulScalar(s float32, a Vec3) Vec3 {
	return Vec3{a.X * s, a.Y * s, a.Z * s}
}

// Mul computes the elementwise product between two vectors
func Mul(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

// Vec3 represents a 3D vector
type Vec3 struct {
	X float32
	Y float32
	Z float32
}

// Pi is the circle constant
const Pi = float32(math.Pi)

func (v Vec3) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}

// RGBA interpretation of the vector
func (v Vec3) RGBA() color.Color {
	return color.RGBA{uint8(v.X * 255), uint8(v.Y * 255), uint8(v.Z * 255), 255}
}

// SquaredLength of the vector
func (v Vec3) SquaredLength() float32 {
	return Dot(v, v)
}

// Length of the vector
func (v Vec3) Length() float32 {
	return Sqrt(v.SquaredLength())
}

// Dot product between two vectors
func Dot(a Vec3, b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross product between two vectors
func Cross(a Vec3, b Vec3) Vec3 {
	return Vec3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

// Normalize a vector
func Normalize(a Vec3) Vec3 {
	return DivScalar(a.Length(), a)
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

// DivScalar computes a / s
func DivScalar(s float32, a Vec3) Vec3 {
	return Vec3{a.X / s, a.Y / s, a.Z / s}
}

// Mul computes the elementwise product between two vectors
func Mul(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// RandomUniform sample a random number in [-1, 1)
func RandomUniform(rng *rand.Rand) float32 {
	return rng.Float32()*2 - 1
}

// RandomPointInUnitSphere samples a random point inside the unit sphere
func RandomPointInUnitSphere(rng *rand.Rand) Vec3 {
	for {
		v := Vec3{RandomUniform(rng), RandomUniform(rng), RandomUniform(rng)}
		if v.SquaredLength() <= 1 {
			return v
		}
	}
}

// Sqrt computes the sqrt of a float32
func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// Deg2Rad converts an angle in degrees to radians
func Deg2Rad(angle float32) float32 {
	return angle / 180 * Pi
}

// Rad2Deg converts an angle in radians to degrees
func Rad2Deg(angle float32) float32 {
	return angle * 180 / Pi
}

// Abs computes the absolute value of a float32
func Abs(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}

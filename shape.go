package main

import (
	"math"
)

// Material of an object
type Material struct {
	Color    Vec3
	Specular float32
}

// Shape in the world
type Shape interface {
	Intersect(Ray) *Hit
}

// Sphere in 3D space
type Sphere struct {
	Position Vec3
	Radius   float32
	Material Material
}

// Intersect check whether the ray intersects the sphere
func (sphere Sphere) Intersect(ray Ray) *Hit {
	a := Dot(ray.Direction, ray.Direction)
	relPos := Sub(ray.Origin, sphere.Position)
	b := 2 * Dot(ray.Direction, relPos)
	// c := -2*mat.Dot(ray.Origin, sphere.Position)*mat.Dot(sphere.Position, sphere.Position) - sphere.Radius*sphere.Radius
	c := Dot(relPos, relPos) - sphere.Radius*sphere.Radius

	discriminant := float64(b*b - 4*a*c)
	if discriminant < 0 {
		return nil
	}

	t := (-b - float32(math.Sqrt(discriminant))) / (2 * a)
	if t < 0 {
		t = (-b + float32(math.Sqrt(discriminant))) / (2 * a)
	}
	if t <= 1e-3 {
		return nil
	}
	normal := Normalize(Sub(ray.At(t), sphere.Position))
	return &Hit{t, normal, sphere.Material}
}

// Plane in 3D space
type Plane struct {
	Normal   Vec3
	Along    float32
	Material Material
}

// Intersect checks if a ray intersects with the plane
func (plane Plane) Intersect(ray Ray) *Hit {
	denom := Dot(Sub(ray.Origin, plane.Normal), ray.Direction)
	if denom <= 0 {
		return nil
	}
	t := 1 / denom
	return &Hit{t, plane.Normal, plane.Material}
}

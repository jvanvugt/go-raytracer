package main

import (
	"math"
	"math/rand"
)

// Material of an object
type Material interface {
	// TODO: put result in struct?
	Scatter(Ray, Hit, *rand.Rand) (didScatter bool, attenuation Vec3, scattered Ray)
}

// Lambertian material
type Lambertian struct {
	Albedo Vec3
}

// Scatter a ray on a lambertian material
func (mat Lambertian) Scatter(ray Ray, hit Hit, rng *rand.Rand) (didScatter bool, attenuation Vec3, scattered Ray) {
	direction := Normalize(Add(RandomPointInUnitSphere(rng), hit.Normal))
	bouncingRay := Ray{hit.Position, direction}
	return true, mat.Albedo, bouncingRay
}

// Metal material
type Metal struct {
	Albedo Vec3
	fuzz   float32
}

// Scatter a ray on a metal material
func (mat Metal) Scatter(ray Ray, hit Hit, rng *rand.Rand) (didScatter bool, attenuation Vec3, scattered Ray) {
	direction := Sub(ray.Direction, MulScalar(2.0*Dot(hit.Normal, ray.Direction), hit.Normal))
	direction = Add(direction, MulScalar(mat.fuzz, RandomPointInUnitSphere(rng)))
	bouncingRay := Ray{hit.Position, direction}
	return Dot(direction, hit.Normal) > 0, mat.Albedo, bouncingRay
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
	return NewHit(t, ray, normal, sphere.Material)
}

// Plane in 3D space
type Plane struct {
	Normal   Vec3
	Along    float32
	Material Material
}

// Intersect checks if a ray intersects with the plane
func (plane Plane) Intersect(ray Ray) *Hit {
	denom := Dot(plane.Normal, ray.Direction)
	if math.Abs(float64(denom)) < 1e-6 {
		return nil
	}
	planePoint := MulScalar(plane.Along, plane.Normal)
	t := (Dot(planePoint, plane.Normal) - Dot(plane.Normal, ray.Origin)) / denom
	if t < 1e-3 {
		return nil
	}
	return NewHit(t, ray, plane.Normal, plane.Material)
}

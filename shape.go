package main

import (
	"log"
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

func reflect(incoming Vec3, normal Vec3) Vec3 {
	return Sub(incoming, MulScalar(2.0*Dot(normal, incoming), normal))
}

// Scatter a ray on a metal material
func (mat Metal) Scatter(ray Ray, hit Hit, rng *rand.Rand) (didScatter bool, attenuation Vec3, scattered Ray) {
	direction := reflect(ray.Direction, hit.Normal)
	direction = Normalize(Add(direction, MulScalar(mat.fuzz, RandomPointInUnitSphere(rng))))
	bouncingRay := Ray{hit.Position, direction}
	return Dot(direction, hit.Normal) > 0, mat.Albedo, bouncingRay
}

// Dielectric materials both reflect and refrect
type Dielectric struct {
	ReflectionIndex float32
}

func refract(incoming Vec3, normal Vec3, niOverNt float32) (didRefract bool, refraction Vec3) {
	dt := Dot(incoming, normal)
	discriminant := 1 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		refraction = Normalize(Sub(MulScalar(niOverNt, Sub(incoming, MulScalar(dt, normal))), MulScalar(Sqrt(discriminant), normal)))
		return true, refraction
	}
	return false, Vec3{}
}

func schlick(cosine float32, reflectionIndex float32) float32 {
	r0 := (1 - reflectionIndex) / (1 + reflectionIndex)
	r0 *= r0
	base := 1 - cosine
	return r0 + (1-r0)*base*base*base*base*base
}

func assertUnitLength(v Vec3) {
	if math.Abs(1-float64(v.Length())) > 1e-3 {
		log.Fatal("Vector is not unit length: ", v.Length())
	}
}

// Scatter a ray on a dielectric
func (mat Dielectric) Scatter(ray Ray, hit Hit, rng *rand.Rand) (didScatter bool, attenuation Vec3, scattered Ray) {
	var outwardNormal Vec3
	var niOverNt float32
	var cosine float32
	if Dot(ray.Direction, hit.Normal) > 0 {
		outwardNormal = MulScalar(-1, hit.Normal)
		niOverNt = mat.ReflectionIndex
		cosine = mat.ReflectionIndex * Dot(ray.Direction, hit.Normal)
	} else {
		outwardNormal = hit.Normal
		niOverNt = 1.0 / mat.ReflectionIndex
		cosine = -Dot(ray.Direction, hit.Normal)
	}

	didRefract, refracted := refract(ray.Direction, outwardNormal, niOverNt)
	if didRefract && schlick(cosine, mat.ReflectionIndex) < rng.Float32() {
		return true, Vec3{1, 1, 1}, Ray{hit.Position, refracted}
	}

	reflected := reflect(ray.Direction, hit.Normal)
	return true, Vec3{1, 1, 1}, Ray{hit.Position, reflected}
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

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return nil
	}

	t := (-b - Sqrt(discriminant)) / (2 * a)
	if t < 0 {
		t = (-b + Sqrt(discriminant)) / (2 * a)
	}
	if t <= 1e-3 {
		return nil
	}

	normal := Normalize(DivScalar(sphere.Radius, Sub(ray.At(t), sphere.Position)))
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

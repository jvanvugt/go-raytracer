package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
)

const imageWidth = 1280
const imageHeight = 720
const aspectRatio = float32(imageWidth) / float32(imageHeight)
const numSamples = 16

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

// Ray from origin in a direction
type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// At computes the point on the ray at t
func (ray *Ray) At(t float32) Vec3 {
	return Add(ray.Origin, MulScalar(t, ray.Direction))
}

// Hit represents data about a ray hitting an object
type Hit struct {
	T        float32
	Normal   Vec3
	Material Material
}

var world = []Shape{
	Sphere{Vec3{1, 1, 3}, 0.5, Material{Vec3{1, 0, 0}, 0.5}},
	Plane{Vec3{0, 1, 0}, -1, Material{Vec3{0, 0, 1}, 0}},
	Sphere{Vec3{0, -1, 2}, 0.5, Material{Vec3{0, 1, 0}, 0}},
	Sphere{Vec3{-3, 2, 2}, 0.5, Material{Vec3{1, 1, 0}, 0}},
	Sphere{Vec3{0, 1, 2}, 0.5, Material{Vec3{1, 0, 1}, 0}},
}

func castRay(ray Ray, bounced int) Vec3 {
	closest := float32(math.MaxFloat32)
	var closestHit *Hit
	for _, shape := range world {
		hit := shape.Intersect(ray)
		if hit != nil && hit.T < closest {
			closest = hit.T
			closestHit = hit
		}
	}
	if closestHit != nil {
		specular := closestHit.Material.Specular
		if specular > 0 && bounced < 8 {
			direction := Add(MulScalar(2*Dot(closestHit.Normal, MulScalar(-1, ray.Direction)), closestHit.Normal), ray.Direction)
			bouncingRay := Ray{ray.At(closestHit.T), direction}
			return Add(MulScalar(1-specular, closestHit.Material.Color),
				MulScalar(specular, castRay(bouncingRay, bounced+1)))
		}
		return closestHit.Material.Color

	} else {
		return Vec3{0.8, 0.8, 1}
	}
}

func getColor(x int, y int, rng *rand.Rand) Vec3 {
	color := Vec3{0, 0, 0}
	for i := 0; i < numSamples; i++ {
		ySample := float32(y) + rng.Float32() - 0.5
		xSample := float32(x) + rng.Float32() - 0.5
		scaledY := -(ySample/float32(imageHeight)*2 - 1)
		scaledX := xSample/float32(imageWidth)*2 - 1
		scaledX *= aspectRatio

		direction := Normalize(Vec3{scaledX, scaledY, 1})
		cameraOrigin := Vec3{0, 0, 0}
		ray := Ray{cameraOrigin, direction}
		color = Add(color, castRay(ray, 0))

	}
	return MulScalar(1/float32(numSamples), color)
}

func processTile(img *image.NRGBA, fromX int, fromY int, toX int, toY int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	rng := rand.New(rand.NewSource(0))
	for y := fromY; y < toY; y++ {
		for x := fromX; x < toX; x++ {
			color := getColor(x, y, rng)
			img.Set(x, y, color.RGBA())
		}
	}
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	go processTile(img, 0, 0, imageWidth/2, imageHeight/2, &waitGroup)
	go processTile(img, imageWidth/2, 0, imageWidth, imageHeight/2, &waitGroup)
	go processTile(img, 0, imageHeight/2, imageWidth/2, imageHeight, &waitGroup)
	go processTile(img, imageWidth/2, imageHeight/2, imageWidth, imageHeight, &waitGroup)
	waitGroup.Wait()
	fmt.Println("Hello world")

	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}

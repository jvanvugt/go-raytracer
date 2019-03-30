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
const fieldOfView float32 = 90.0
const numSamples = 100
const maxBounces = 50

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
	Position Vec3
	Normal   Vec3
	Material Material
}

// NewHit creates a Hit object
func NewHit(t float32, ray Ray, normal Vec3, material Material) *Hit {
	return &Hit{
		t,
		ray.At(t),
		normal,
		material,
	}
}

// Camera to shoot rays from
type Camera struct {
	Position   Vec3
	BottomLeft Vec3
	PixelStepX Vec3
	PixelStepY Vec3
}

func setupCamera(cameraPos Vec3, cameraTarget Vec3, up Vec3) Camera {
	cameraDirection := Normalize(Sub(cameraTarget, cameraPos))
	horizontalDirection := Cross(Normalize(up), cameraDirection)
	verticalDirection := Cross(Normalize(cameraDirection), Normalize(horizontalDirection))
	halfWidth := float32(math.Tan(float64(Deg2Rad(fieldOfView)) / 2.0))
	halfHeight := halfWidth * float32(imageHeight) / float32(imageWidth)
	pixelStepX := MulScalar(2*halfWidth/(imageWidth-1), horizontalDirection)
	pixelStepY := MulScalar(2*halfHeight/(imageHeight-1), verticalDirection)
	bottomLeft := Sub(Sub(cameraDirection, MulScalar(halfWidth, horizontalDirection)), MulScalar(halfHeight, verticalDirection))
	return Camera{cameraPos, bottomLeft, pixelStepX, pixelStepY}
}

func (camera *Camera) getRay(x float32, y float32) Ray {
	direction := Add(Add(camera.BottomLeft, MulScalar(x, camera.PixelStepX)), MulScalar(y, camera.PixelStepY))
	return Ray{camera.Position, Normalize(direction)}
}

// My own scene
// var world = []Shape{
// 	Sphere{Vec3{1, 1, 3}, 0.5, Metal{Vec3{1, 1, 1}, 0.3}},
// 	Plane{Vec3{0, 1, 0}, -1, Lambertian{Vec3{0.7, 0.8, 1.0}}},
// 	Sphere{Vec3{0, -0.5, 2}, 0.5, Lambertian{Vec3{0, 1, 0}}},
// 	Sphere{Vec3{-3, 2, 2}, 0.5, Lambertian{Vec3{1, 1, 0}}},
// 	Sphere{Vec3{0, 1, 2}, 0.5, Lambertian{Vec3{1, 0, 1}}},
// }

// Two metal balls
// var world = []Shape{
// 	Plane{Vec3{0, 1, 0}, -1, Lambertian{Vec3{140 / 255., 245 / 255., 98 / 255.}}},
// 	Sphere{Vec3{-2, 0, 2}, 1, Metal{Vec3{1, 1, 1}, 0.2}},
// 	Sphere{Vec3{0, 0, 2}, 1, Lambertian{Vec3{255 / 255., 200 / 255., 210 / 255.}}},
// 	Sphere{Vec3{2, 0, 2}, 1, Metal{Vec3{0.8, 0.75, 1}, 0}},
// }

// Dielectrics
var world = []Shape{
	Sphere{Vec3{0, 0, 2}, 0.5, Lambertian{Vec3{0.1, 0.2, 0.5}}},
	Sphere{Vec3{0, -100.5, 1}, 100, Lambertian{Vec3{0.8, 0.8, 0.0}}},
	Sphere{Vec3{1, 0, 2}, 0.5, Metal{Vec3{0.8, 0.6, 0.2}, 0}},
	Sphere{Vec3{-1, 0, 2}, 0.45, Dielectric{1.5}},
}

func castRay(ray Ray, rng *rand.Rand, bounced int) Vec3 {
	if bounced > maxBounces {
		return Vec3{0, 0, 0}
	}
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
		didScatter, attenuation, scatteredRay := closestHit.Material.Scatter(ray, *closestHit, rng)
		if didScatter {
			return Mul(attenuation, castRay(scatteredRay, rng, bounced+1))
		}
		return Vec3{0, 0, 0}
	}

	return Add(MulScalar((ray.Direction.Y+1)/2, Vec3{0.6, 0.6, 1}), MulScalar(1-(ray.Direction.Y+1)/2, Vec3{1, 1, 1}))
}

func getColor(camera *Camera, x int, y int, rng *rand.Rand) Vec3 {
	color := Vec3{0, 0, 0}
	for i := 0; i < numSamples; i++ {
		ray := camera.getRay(float32(x)+rng.Float32()-0.5, float32(y)+rng.Float32()-0.5)
		color = Add(color, castRay(ray, rng, 0))

	}
	return DivScalar(float32(numSamples), color)
}

func processTile(img *image.NRGBA, camera *Camera, fromX int, fromY int, toX int, toY int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	rng := rand.New(rand.NewSource(0))
	for y := fromY; y < toY; y++ {
		for x := fromX; x < toX; x++ {
			color := getColor(camera, x, y, rng)
			gammaCorrectedColor := Vec3{Sqrt(color.X), Sqrt(color.Y), Sqrt(color.Z)}
			img.Set(x, imageHeight-y-1, gammaCorrectedColor.RGBA())
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

	cameraPos := Vec3{0, 0, 0}
	target := Vec3{0, 0, 1}
	up := Vec3{0, 1, 0}
	camera := setupCamera(cameraPos, target, up)

	img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	go processTile(img, &camera, 0, 0, imageWidth/2, imageHeight/2, &waitGroup)
	go processTile(img, &camera, imageWidth/2, 0, imageWidth, imageHeight/2, &waitGroup)
	go processTile(img, &camera, 0, imageHeight/2, imageWidth/2, imageHeight, &waitGroup)
	go processTile(img, &camera, imageWidth/2, imageHeight/2, imageWidth, imageHeight, &waitGroup)
	waitGroup.Wait()
	fmt.Println("Hello world")

	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, img)
}

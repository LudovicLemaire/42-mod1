package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func GetDistance(x, y, z float64) float64 {
	return math.Sqrt(math.Pow(x-0.0, 2) + math.Pow(y-0, 2) + math.Pow(z-0, 2))
}

func AddCube(pos Vec3u32, rgb mgl32.Vec3, points []float32) []float32 {
	rgb = generateVarianteColorBillowNoise(rgb, float64(pos[0]), float64(pos[1]), float64(pos[2]))
	newCube := []float32{
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 1.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
	}
	points = append(points, newCube...)
	return points
}

func AddPlane(pos Vec3u32, rgb mgl32.Vec3, points []float32) []float32 {
	newCube := []float32{
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],

		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 0.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		1.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
		0.0 + float32(pos[0]), 0.0 + float32(pos[1]), 1.0 + float32(pos[2]), rgb[0], rgb[1], rgb[2],
	}
	points = append(points, newCube...)
	return points
}

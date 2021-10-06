package main

import (
	"math/rand"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func AllScenarios() {
	if rainScenario {
		points_water = []float32{}
		for x := 0; x < simulationSize; x++ {
			for z := 0; z < simulationSize; z++ {
				if rand.Float64() > 0.999 {
					if !groundMap[Vec3i32{int32(x), int32(simulationSize + 1), int32(z)}] {
						waterMap[Vec3i32{int32(x), int32(simulationSize + 1), int32(z)}] = true
					}
				}
			}
		}
		for key := range waterMap {
			points_water = AddCube(key, cWater, points_water)
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
	}

	if waveScenario {
		points_water = []float32{}
		for x := 0; x < simulationSize; x++ {
			for z := 0; z < 3; z++ {
				for y := 0; y < simulationSize; y++ {
					if !groundMap[Vec3i32{int32(x), int32(y), int32(z)}] {
						waterMap[Vec3i32{int32(x), int32(y), int32(z)}] = true
					}
				}
			}
		}
		for key := range waterMap {
			points_water = AddCube(key, cWater, points_water)
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
	}
	if floodScenario {
		elapsedFloodScenario := time.Since(startFloodScenario)
		var floodLevel int = int(elapsedFloodScenario.Seconds() / 1)
		if currFloodLevel != floodLevel && currFloodLevel < simulationSize {
			currFloodLevel = floodLevel
			points_water = []float32{}
			for x := 0; x < simulationSize; x++ {
				for z := 0; z < simulationSize; z++ {
					if !groundMap[Vec3i32{int32(x), int32(currFloodLevel), int32(z)}] {
						waterMap[Vec3i32{int32(x), int32(currFloodLevel), int32(z)}] = true
					}
				}

			}
			for key := range waterMap {
				points_water = AddCube(key, cWater, points_water)
			}
			gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
			gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
		}
	}
}

func SpawnerScenarios(keys Keys) {
	if keys.right == "hold" || keys.left == "hold" || keys.up == "hold" || keys.down == "hold" || keys.add == "hold" || keys.minus == "hold" || keys.multiply == "hold" || keys.divide == "hold" {
		points_waterSpawner = []float32{}
		waterSpawnerMap = make(map[Vec3i32]bool)
		waterSpawnerMap[Vec3i32{0, -50000, 0}] = true
		// generate new water spawner positions
		for x := 0; x < sizeWs; x++ {
			for z := 0; z < sizeWs; z++ {
				if (x+z)%2 == 0 {
					if z+zOffsetWS < simulationSize && x+xOffsetWS < simulationSize && z+zOffsetWS >= 0 && x+xOffsetWS >= 0 {
						waterSpawnerMap[Vec3i32{int32(x + xOffsetWS), int32(yOffsetWs), int32(z + zOffsetWS)}] = true
					}
				}
			}
		}
		for key := range waterSpawnerMap {
			points_waterSpawner = AddPlane(key, cPlaneWaterSpawner, points_waterSpawner)
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_waterSpawner)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_waterSpawner), gl.Ptr(points_waterSpawner), gl.DYNAMIC_DRAW)
	}
	if keys.f == "hold" {
		// add Water from Spawner
		points_water = []float32{}
		for x := 0; x < sizeWs; x++ {
			for z := 0; z < sizeWs; z++ {
				if z+zOffsetWS < simulationSize && x+xOffsetWS < simulationSize && z+zOffsetWS >= 0 && x+xOffsetWS >= 0 {
					if !groundMap[Vec3i32{int32(x + xOffsetWS), int32(yOffsetWs), int32(z + zOffsetWS)}] {
						if rand.Float64() > 0.75 {
							waterMap[Vec3i32{int32(x + xOffsetWS), int32(yOffsetWs), int32(z + zOffsetWS)}] = true
						}
					}
				}
			}
		}

		for key := range waterMap {
			points_water = AddCube(key, cWater, points_water)
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
	}

	if keys.g == "hold" {
		// remove Water from Spawner
		points_water = []float32{}
		for x := 0; x < sizeWs; x++ {
			for z := 0; z < sizeWs; z++ {
				if z+zOffsetWS < simulationSize && x+xOffsetWS < simulationSize && z+zOffsetWS >= 0 && x+xOffsetWS >= 0 {
					if !groundMap[Vec3i32{int32(x + xOffsetWS), int32(yOffsetWs), int32(z + zOffsetWS)}] {
						waterMap[Vec3i32{int32(x + xOffsetWS), int32(yOffsetWs), int32(z + zOffsetWS)}] = false
					}
				}
			}
		}

		for key := range waterMap {
			if waterMap[key] {
				points_water = AddCube(key, cWater, points_water)
			}
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
	}

	if keys.h == "hold" {
		// remove Ground from Spawner
		points_ground = []float32{}
		for x := 0; x < sizeWs; x++ {
			for z := 0; z < sizeWs; z++ {
				if z+zOffsetWS < simulationSize && x+xOffsetWS < simulationSize && z+zOffsetWS >= 0 && x+xOffsetWS >= 0 {
					groundMap[Vec3i32{int32(x + xOffsetWS), int32(yOffsetWs), int32(z + zOffsetWS)}] = false
				}
			}
		}
		for key := range groundMap {
			if groundMap[key] {
				points_ground = AddCube(key, cGround, points_ground)
			}
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_ground)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_ground), gl.Ptr(points_ground), gl.DYNAMIC_DRAW)
	}
}

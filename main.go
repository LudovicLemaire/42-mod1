package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"mod1/glfont"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type ScreenData struct {
	font *glfont.Font
}

const (
	width          = 1600
	height         = 900
	simulationSize = 175
)

type Mod1 struct {
	pos mgl32.Vec3
	rot mgl32.Vec3
}

var (
	points_ground                = []float32{}
	points_water                 = []float32{}
	points_delimiter             = []float32{}
	_oldMousePosX, _oldMousePosY float64
	cGround                      = mgl32.Vec3{0.29, 0.68, 0.31}
	cWater                       = mgl32.Vec3{0.12, 0.58, 0.94}
	cDelimiter                   = mgl32.Vec3{0.95, 0.26, 0.21}
)

type Vec3u32 [3]int32

type ColorRGB struct {
	r float32
	g float32
	b float32
}

type GameValues struct {
	speed       float32
	polygonMode bool
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var screenData ScreenData
	var err error
	NoiseInitPermtables(445)

	groundMap := make(map[Vec3u32]bool)
	waterMap := make(map[Vec3u32]uint8)
	delimiterMap := make(map[Vec3u32]bool)
	//waterMapPassive := make(map[Vec3u32]uint8)

	// Simulation delimiter \\
	for x := 0; x < simulationSize; x++ {
		for z := 0; z < simulationSize; z++ {
			if (x+z)%2 == 0 {
				delimiterMap[Vec3u32{int32(x), int32(0), int32(z)}] = true
				delimiterMap[Vec3u32{int32(x), int32(simulationSize), int32(z)}] = true
			}
		}

	}
	for key := range delimiterMap {
		points_delimiter = AddPlane(key, cDelimiter, points_delimiter)
	}
	// Simulation delimiter //

	/*
		// ground
		for x := 0; x < simulationSize; x++ {
			for z := 0; z < simulationSize; z++ {
				y := Noise2dSimplex(float64(x), float64(z), 0.0, 0.15, 0.15, 3, 0) * simulationSize / 4
				m[Vec3u32{uint32(x), uint32(y), uint32(z)}] = 1

			}
		}
	*/
	// Water \\
	/*
		for x := 0; x < simulationSize; x++ {
			for z := 0; z < simulationSize; z++ {
				y := Noise2dSimplex(float64(x), float64(z), 0.0, 0.15, 0.15, 3, 0) * simulationSize * 1.5
				waterMap[Vec3u32{int32(x), int32(y), int32(z)}] = 1
			}
		}
	*/
	waterMap[Vec3u32{0, -50000, 0}] = 1
	for key := range waterMap {
		points_water = AddCube(key, cWater, points_water)
	}
	// Water //

	// Ground \\
	/*
		// cup
		for x := 7; x < 14; x++ {
			for z := 7; z < 14; z++ {
				groundMap[Vec3u32{int32(x), int32(3 + 4), int32(z)}] = 1
			}
		}
		for x := 6; x < 15; x++ {
			for z := 6; z < 15; z++ {
				groundMap[Vec3u32{int32(x), int32(4 + 4), int32(z)}] = 1
			}
		}
		for x := 7; x < 14; x++ {
			for z := 7; z < 14; z++ {
				groundMap[Vec3u32{int32(x), int32(4 + 4), int32(z)}] = 0
			}
		}
	*/

	for x := 0; x < simulationSize; x++ {
		for z := 0; z < simulationSize; z++ {
			y := Noise2dSimplex(float64(x), float64(z), 0, 1.0, 0.0155, 5, 2) * simulationSize / 4
			for i := y; i >= 0; i-- {
				groundMap[Vec3u32{int32(x), int32(i), int32(z)}] = true
			}

		}
	}
	groundMap[Vec3u32{0, -50000, 0}] = true
	for key := range groundMap {
		if groundMap[key] {
			points_ground = AddCube(key, cGround, points_ground)
		}
	}
	// Ground //

	// Check if level is filled \\
	// filledMapBase := make(map[int64]int64)
	// for key := range groundMap {
	// 	if groundMap[key] {
	// 		filledMapBase[int64(key[1])]++
	// 	}
	// }
	// for key := range filledMapBase {
	// 	fmt.Println(key, ": ", filledMapBase[key])
	// }
	// Check if level is filled //

	var Mod1 Mod1
	var colorRColorRGB ColorRGB
	colorRColorRGB.r, colorRColorRGB.g, colorRColorRGB.b = 0.0, 0.0, 0.0
	Mod1.pos[0], Mod1.pos[1], Mod1.pos[2] = simulationSize/2, -simulationSize/2, -simulationSize/2
	Mod1.rot[1] = 1.5

	window := InitGlfw()
	defer glfw.Terminate()
	program := InitOpenGL()

	whiteUniform := gl.GetUniformLocation(program, gl.Str("white\x00"))

	cameraId := gl.GetUniformLocation(program, gl.Str("camera\x00"))

	gl.PointSize(5)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	var vbo_ground uint32
	gl.GenBuffers(1, &vbo_ground)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_ground)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_ground), gl.Ptr(points_ground), gl.STATIC_DRAW)

	var vbo_water uint32
	gl.GenBuffers(1, &vbo_water)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.STATIC_DRAW)

	var vbo_delimiter uint32
	gl.GenBuffers(1, &vbo_delimiter)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_delimiter)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_delimiter), gl.Ptr(points_delimiter), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Disable(gl.CULL_FACE)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	screenData.font, err = glfont.LoadFont("Assets/Fonts/SourceCodePro-Regular.ttf", int32(14), int(width), int(height))
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	var keys Keys
	initKeys(&keys)
	var gameValues GameValues
	initGameValues(&gameValues)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		EventsKeyboard(&Mod1, &colorRColorRGB, &keys, &gameValues, whiteUniform)
		EventsMouse(&Mod1)
		setCamera(cameraId, &Mod1)

		if keys.v == "active" || keys.b == "hold" {
			// do iteration for water
			//Total := time.Now()
			//startCreateMap := time.Now()
			points_water = []float32{}
			waterMapNew := make(map[Vec3u32]uint8)
			//elapsedCreateMap := time.Since(startCreateMap)
			//fmt.Println("Elapsed CreateMap: ", elapsedCreateMap)
			//startMoveAllWater := time.Now()
			moveWater(waterMap, waterMapNew, groundMap)
			//elapsedMoveAllWater := time.Since(startMoveAllWater)
			//fmt.Println("Elapsed MoveAllWater: ", elapsedMoveAllWater)
			//startCopyWater := time.Now()
			waterMap = waterMapNew

			//for k, v := range waterMapPassive {
			//	waterMapNew[k] = v
			//}

			//elapsedCopyWater := time.Since(startCopyWater)
			//fmt.Println("Elapsed CopyWater: ", elapsedCopyWater)
			//startCreateCube := time.Now()
			for key := range waterMapNew {
				if waterMapNew[key] == 1 {
					points_water = AddCube(key, cWater, points_water)
				}
			}
			//elapsedCreateCube := time.Since(startCreateCube)
			//fmt.Println("Elapsed CreateCube: ", elapsedCreateCube)
			//startBufferData := time.Now()
			gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
			gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
			//elapsedBufferData := time.Since(startBufferData)
			//fmt.Println("Elapsed BufferData: ", elapsedBufferData)
			//elapsedTotal := time.Since(startTotal)
			//fmt.Println("Elapsed Total: ", elapsedTotal)
			//fmt.Println()
		}

		if keys.c == "active" {
			// add water
			points_water = []float32{}

			for x := 0; x < simulationSize; x++ {
				for z := 0; z < simulationSize; z++ {
					y := Noise2dSimplex(float64(x), float64(z), rand.Float64()/4, 0.15, 0.15, 3, 0) * simulationSize * 1.75
					if rand.Float64() > 0.75 {
						waterMap[Vec3u32{int32(x), int32(y), int32(z)}] = 1
					}
				}
			}
			for key := range waterMap {
				points_water = AddCube(key, cWater, points_water)
			}

			gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
			gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
		}

		if keys.x == "hold" {
			// add water
			points_water = []float32{}

			//waterMap[Vec3u32{0, simulationSize, 0}] = 1
			//waterMap[Vec3u32{1, simulationSize, 0}] = 1
			//waterMap[Vec3u32{0, simulationSize, 1}] = 1
			//waterMap[Vec3u32{1, simulationSize, 1}] = 1
			//waterMap[Vec3u32{2, simulationSize, 2}] = 1

			for i := 0; i < simulationSize; i++ {
				for y := 0; y < simulationSize/3; y++ {
					waterMap[Vec3u32{int32(i), int32(y), 0}] = 1
				}
			}

			for key := range waterMap {
				points_water = AddCube(key, cWater, points_water)
			}

			gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
			gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
		}
		if keys.z == "hold" {
			// add water
			points_water = []float32{}

			wow := simulationSize / 10
			for i := (simulationSize / 2) - wow; i < (simulationSize/2)+wow; i++ {
				for ip := (simulationSize / 2) - wow; ip < (simulationSize/2)+wow; ip++ {
					if rand.Float64() > 0.95 {
						waterMap[Vec3u32{int32(i), simulationSize, int32(ip)}] = 1
					}
				}
			}

			for key := range waterMap {
				points_water = AddCube(key, cWater, points_water)
			}

			gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
			gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
		}

		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_ground)
		gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
		gl.EnableVertexAttribArray(1)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points_ground)/3))

		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_delimiter)
		gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
		gl.EnableVertexAttribArray(1)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points_delimiter)/3))

		gl.BindVertexArray(vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
		gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
		gl.EnableVertexAttribArray(1)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points_water)/3))

		// print text
		gl.Finish()
		if keys.graveAccent == "hold" {
			displayScreenInfo(screenData, len(waterMap), len(groundMap))
		}

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func setCamera(cameraId int32, Mod1 *Mod1) {
	camera := mgl32.HomogRotate3D(Mod1.rot.X(), mgl32.Vec3{1, 0, 0})
	camera = camera.Mul4(mgl32.HomogRotate3D(Mod1.rot.Y(), mgl32.Vec3{0, 1, 0}))
	camera = camera.Mul4(mgl32.Translate3D(Mod1.pos.X(), Mod1.pos.Y(), Mod1.pos.Z()))

	projection := mgl32.Perspective(mgl32.DegToRad(80.0), float32(width)/float32(height), 0.1, 50000)
	view := projection.Mul4(camera)

	gl.UniformMatrix4fv(cameraId, 1, false, &view[0])
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShaderSource := _getShaderSource("Shaders/qd.vertex.glsl")
	fragmentShaderSource := _getShaderSource("Shaders/qd.fragment.glsl")

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// initGlfw initializes glfw and returns a Window to use.
func InitGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "woop woop ft_Mod1", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func init() {
	runtime.LockOSThread()
}

/*
	fmt.Println(len(points) / 6 / 3)
*/

func displayScreenInfo(screenData ScreenData, totalWater, totalGround int) {
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	screenData.font.SetColor(0.29, 0.68, 0.31, 1)
	screenData.font.Printf(6, 20, 1.0, "Ground:")
	screenData.font.SetColor(1, 1, 1, 1)
	screenData.font.Printf(76, 20, 1.0, "%d", totalGround)
	screenData.font.SetColor(0.12, 0.58, 0.94, 1)
	screenData.font.Printf(6, 40, 1.0, "Water:")
	screenData.font.SetColor(1, 1, 1, 1)
	screenData.font.Printf(76, 40, 1.0, "%d", totalWater)
}

/*
start := time.Now()
elapsed := time.Since(start)
*/

func moveWater(waterMap, waterMapNew map[Vec3u32]uint8, groundMap map[Vec3u32]bool) {
	for key := range waterMap {
		if waterMap[key] == 1 {
			if waterMap[Vec3u32{key[0], key[1] - 1, key[2]}] == 0 && key[1]-1 >= 0 &&
				!groundMap[Vec3u32{key[0], key[1] - 1, key[2]}] &&
				waterMapNew[Vec3u32{key[0], key[1] - 1, key[2]}] == 0 {

				waterMapNew[Vec3u32{key[0], key[1] - 1, key[2]}] = 1
				waterMap[key] = 0
			} else if groundMap[Vec3u32{key[0], key[1] - 1, key[2]}] &&
				groundMap[Vec3u32{key[0], key[1] - 1, key[2] + 1}] &&
				groundMap[Vec3u32{key[0], key[1] - 1, key[2] - 1}] &&
				groundMap[Vec3u32{key[0] + 1, key[1] - 1, key[2]}] &&
				groundMap[Vec3u32{key[0] - 1, key[1] - 1, key[2]}] &&

				groundMap[Vec3u32{key[0] - 1, key[1] - 1, key[2] + 1}] &&
				groundMap[Vec3u32{key[0] - 1, key[1] - 1, key[2] - 1}] &&
				groundMap[Vec3u32{key[0] + 1, key[1] - 1, key[2] + 1}] &&
				groundMap[Vec3u32{key[0] + 1, key[1] - 1, key[2] - 1}] {
				waterMapNew[key] = 1
			} else {
				visitedMap := make(map[Vec3u32]bool)
				var nextPosArray [4]Vec3u32 = [4]Vec3u32{
					{key[0] + 1, key[1], key[2]},
					{key[0] - 1, key[1], key[2]},
					{key[0], key[1], key[2] + 1},
					{key[0], key[1], key[2] - 1},
				}
				rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })
				didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[0], 0)
				if didFound {
					if waterMapNew[posFound] == 0 && waterMap[posFound] == 0 && !groundMap[posFound] {
						waterMapNew[posFound] = 1
						waterMap[key] = 0
					} else {
						waterMapNew[key] = 1
					}
				} else {
					didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[1], 0)
					if didFound {
						if waterMapNew[posFound] == 0 && waterMap[posFound] == 0 && !groundMap[posFound] {
							waterMapNew[posFound] = 1
							waterMap[key] = 0
						} else {
							waterMapNew[key] = 1
						}
					} else {
						didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[2], 0)
						if didFound {
							if waterMapNew[posFound] == 0 && waterMap[posFound] == 0 && !groundMap[posFound] {
								waterMapNew[posFound] = 1
								waterMap[key] = 0
							} else {
								waterMapNew[key] = 1
							}
						} else {
							didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, key, key, nextPosArray[3], 0)
							if didFound {
								if waterMapNew[posFound] == 0 && waterMap[posFound] == 0 && !groundMap[posFound] {
									waterMapNew[posFound] = 1
									waterMap[key] = 0
								} else {
									waterMapNew[key] = 1
								}
							} else {
								waterMapNew[key] = 1
							}
						}
					}
				}
			}
		}
	}
}

func search(waterMap, waterMapNew map[Vec3u32]uint8, groundMap, visitedMap map[Vec3u32]bool, initialPos, oldPos, currPos Vec3u32, nbIteration int) (bool, Vec3u32) {
	var maxIteration int = 100
	if visitedMap[currPos] || groundMap[currPos] || waterMap[currPos] == 1 || waterMapNew[currPos] == 1 ||
		currPos[0] < 0 || currPos[1] < 0 || currPos[2] < 0 ||
		currPos[0] >= simulationSize || currPos[1] >= simulationSize || currPos[2] >= simulationSize ||
		nbIteration > maxIteration {
		return false, Vec3u32{0, 0, 0}
	} else {
		visitedMap[currPos] = true
	}

	//if groundMap[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] || waterMap[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] == 1 || waterMapNew[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] == 1 {
	//	return true, Vec3u32{currPos[0], currPos[1], currPos[2]}
	//}

	if waterMap[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] == 0 &&
		!groundMap[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] &&
		waterMapNew[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] == 0 &&
		currPos[1]-1 >= 0 {

		if waterMapNew[Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}] == 0 {
			return true, Vec3u32{currPos[0], currPos[1] - 1, currPos[2]}
		} else {
			fmt.Println("\033[1;31mUnintended behavior:\033[0m code 1")
			return false, Vec3u32{0, 0, 0}
		}
	} else {
		var nextPosArray [4]Vec3u32 = [4]Vec3u32{
			{currPos[0] + 1, currPos[1], currPos[2]},
			{currPos[0] - 1, currPos[1], currPos[2]},
			{currPos[0], currPos[1], currPos[2] + 1},
			{currPos[0], currPos[1], currPos[2] - 1},
		}
		rand.Shuffle(len(nextPosArray), func(i, j int) { nextPosArray[i], nextPosArray[j] = nextPosArray[j], nextPosArray[i] })

		newPos := nextPosArray[0]
		didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1)
		if didFound {
			return true, Vec3u32{posFound[0], posFound[1], posFound[2]}
		} else {
			newPos := nextPosArray[1]
			didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1)
			if didFound {
				return true, Vec3u32{posFound[0], posFound[1], posFound[2]}
			} else {
				newPos := nextPosArray[2]
				didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1)
				if didFound {
					return true, Vec3u32{posFound[0], posFound[1], posFound[2]}
				} else {
					newPos := nextPosArray[3]
					didFound, posFound := search(waterMap, waterMapNew, groundMap, visitedMap, initialPos, currPos, newPos, nbIteration+1)
					if didFound {
						return true, Vec3u32{posFound[0], posFound[1], posFound[2]}
					} else {
						return false, Vec3u32{0, 0, 0}
					}
				}
			}
		}
	}
}

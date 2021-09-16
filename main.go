package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width          = 800
	height         = 500
	simulationSize = 250
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

type Vec3u32 [3]uint32

var m map[Vec3u32]uint8
var mMutex = sync.RWMutex{}

type ColorRGB struct {
	r float32
	g float32
	b float32
}

type GameValues struct {
	speed       float32
	polygonMode bool
}

func findMe() {

}

func tamer(points_water []float32, vbo_water uint32) {
	points_water = []float32{}
	mCp := make(map[Vec3u32]uint8, len(m))
	for key := range m {
		if m[key] == 2 {
			if (m[Vec3u32{key[0], key[1] - 1, key[2]}] == 0) {
				mCp[Vec3u32{key[0], key[1] - 1, key[2]}] = 2
			} else {
				mCp[key] = m[key]
			}
		} else {
			mCp[key] = m[key]
		}
	}
	m = mCp
	//
	for key := range m {
		if m[key] == 2 {
			points_water = AddCube(key, cWater, points_water)
		}
	}
	//elapsed := time.Since(start)
	//log.Printf("time: %s", elapsed)
	// 110ms
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_water)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points_water), gl.Ptr(points_water), gl.DYNAMIC_DRAW)
	// 8ms
}

func main() {
	rand.Seed(time.Now().UnixNano())
	NoiseInitPermtables(42)

	m = make(map[Vec3u32]uint8)

	// simulation delimiter
	for x := 0; x < simulationSize; x++ {
		for z := 0; z < simulationSize; z++ {
			if (x+z)%2 == 0 {
				m[Vec3u32{uint32(x), uint32(0), uint32(z)}] = 3
				m[Vec3u32{uint32(x), uint32(simulationSize), uint32(z)}] = 3
			}
		}
	}
	// ground
	for x := 0; x < simulationSize; x++ {
		for z := 0; z < simulationSize; z++ {
			y := Noise2dSimplex(float64(x), float64(z), 0.0, 0.15, 0.15, 3, 0) * simulationSize / 4
			m[Vec3u32{uint32(x), uint32(y), uint32(z)}] = 1

		}
	}
	// water
	for x := 0; x < simulationSize; x++ {
		for z := 0; z < simulationSize; z++ {
			y := Noise2dSimplex(float64(x), float64(z), 0.0, 0.15, 0.15, 3, 0) * simulationSize * 1.5
			m[Vec3u32{uint32(x), uint32(y), uint32(z)}] = 2

		}
	}

	for key, value := range m {
		//fmt.Println(key, value)
		if value == 1 {
			points_ground = AddCube(key, cGround, points_ground)
		} else if value == 2 {
			points_water = AddCube(key, cWater, points_water)
		} else if value == 3 {
			points_delimiter = AddPlane(key, cDelimiter, points_delimiter)
		}
	}

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

		if keys.v == "hold" {
			// do iteration for water
			tamer(points_water, vbo_water)
			//fmt.Println(len(points) / 6 / 3)
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

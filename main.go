package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 500
)

var (
	points = []float32{}
	/*
		points = []float32{
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
		}
	*/
	_oldMousePosX, _oldMousePosY float64
)

type Cbo struct {
	pos mgl32.Vec3
	rot mgl32.Vec3
}

type ColorRGB struct {
	r float32
	g float32
	b float32
}

type PlanetData struct {
	seaLevel    float64
	totalPoints int
	scale       float32
	dt          SimplexDt
	dtBillow    SimplexDt
	pos         mgl32.Mat3
	colorGround ColorRGB
	colorSea    ColorRGB
}

type GameValues struct {
	speed       float32
	polygonMode bool
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
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
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "woop woop ft_cbo", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func getDistance(x, y, z float64) float64 {
	return math.Sqrt(math.Pow(x-0.0, 2) + math.Pow(y-0, 2) + math.Pow(z-0, 2))
}

func generatePlanet(points *[]float32, pDt PlanetData) {
	var angleIncrement float64 = 10.1664073846 // TAU * golden ratio

	for i := 0; i <= pDt.totalPoints; i++ {
		var t float64 = float64(i) / float64(pDt.totalPoints)
		var angle1 float64 = math.Acos(1 - 2*float64(t))
		var angle2 float64 = angleIncrement * float64(i)

		var x float32 = float32(math.Sin(angle1) * math.Cos(angle2))
		var y float32 = float32(math.Sin(angle1) * math.Sin(angle2))
		var z float32 = float32(math.Cos(angle1))
		var r float32 = float32(math.Pow(rand.Float64()/2.5, 1.5))
		var g float32 = rand.Float32()
		var b float32 = rand.Float32()

		simplexNoise := (Noise3dSimplex(float64(x+250), float64(y+250), float64(z+250), 0, pDt.dt))
		simplexNoiseBiome := float32(Noise3dSimplex(float64(x+250), float64(y+250), float64(z+250), 1, pDt.dtBillow))
		if simplexNoise > pDt.seaLevel {
			x *= (1 + float32(simplexNoise))
			y *= (1 + float32(simplexNoise))
			z *= (1 + float32(simplexNoise))
		} else {
			x *= (1 + float32(0.3))
			y *= (1 + float32(0.3))
			z *= (1 + float32(0.3))
		}
		rd := rand.Float32() * 7.5

		distanceFromCenter := getDistance(float64(x), float64(y), float64(z))
		if distanceFromCenter < pDt.seaLevel+1 {
			r = (rd+pDt.colorSea.r)/255.0 + simplexNoiseBiome
			g = (rd+pDt.colorSea.g)/255.0 + simplexNoiseBiome
			b = (rd+pDt.colorSea.b)/255.0 + simplexNoiseBiome
		} else {
			r = (rd+pDt.colorGround.r)/255.0 + simplexNoiseBiome
			g = (rd+pDt.colorGround.g)/255.0 + simplexNoiseBiome
			b = (rd+pDt.colorGround.b)/255.0 + simplexNoiseBiome
		}
		*points = append(*points, (x*pDt.scale)+pDt.pos[0], (y*pDt.scale)+pDt.pos[1], (z*pDt.scale)+pDt.pos[2], r, g, b)
	}
}

func subdivideTriangles(points []float32) []float32 {
	var newPoints []float32
	for i := 0; i < len(points); i += (6 * 3) {
		center := mgl32.Vec3{0, 0, 0}
		pt1 := mgl32.Vec3{points[i], points[i+1], points[i+2]}
		pt2 := mgl32.Vec3{points[i+6], points[i+7], points[i+8]}
		pt3 := mgl32.Vec3{points[i+12], points[i+13], points[i+14]}

		pt4 := mgl32.Vec3{(pt1[0] + pt2[0]) / 2, (pt1[1] + pt2[1]) / 2, (pt1[2] + pt2[2]) / 2}
		pt5 := mgl32.Vec3{(pt2[0] + pt3[0]) / 2, (pt2[1] + pt3[1]) / 2, (pt2[2] + pt3[2]) / 2}
		pt6 := mgl32.Vec3{(pt1[0] + pt3[0]) / 2, (pt1[1] + pt3[1]) / 2, (pt1[2] + pt3[2]) / 2}

		v1 := pt1.Sub(center).Normalize()
		pt1 = v1.Mul(7500).Add(center)
		v2 := pt2.Sub(center).Normalize()
		pt2 = v2.Mul(7500).Add(center)
		v3 := pt3.Sub(center).Normalize()
		pt3 = v3.Mul(7500).Add(center)
		v4 := pt4.Sub(center).Normalize()
		pt4 = v4.Mul(7500).Add(center)
		v5 := pt5.Sub(center).Normalize()
		pt5 = v5.Mul(7500).Add(center)
		v6 := pt6.Sub(center).Normalize()
		pt6 = v6.Mul(7500).Add(center)

		simplexNoise1 := float32(Noise3dSimplex(float64(pt1[0]+1000), float64(pt1[1]+25000), float64(pt1[2]+1000), 0, SimplexDt{n: 0.0, a: 0.05, freq: 0.00005, oct: 8}))
		simplexNoise2 := float32(Noise3dSimplex(float64(pt2[0]+1000), float64(pt2[1]+25000), float64(pt2[2]+1000), 0, SimplexDt{n: 0.0, a: 0.05, freq: 0.00005, oct: 8}))
		simplexNoise3 := float32(Noise3dSimplex(float64(pt3[0]+1000), float64(pt3[1]+25000), float64(pt3[2]+1000), 0, SimplexDt{n: 0.0, a: 0.05, freq: 0.00005, oct: 8}))
		simplexNoise4 := float32(Noise3dSimplex(float64(pt4[0]+1000), float64(pt4[1]+25000), float64(pt4[2]+1000), 0, SimplexDt{n: 0.0, a: 0.05, freq: 0.00005, oct: 8}))
		simplexNoise5 := float32(Noise3dSimplex(float64(pt5[0]+1000), float64(pt5[1]+25000), float64(pt5[2]+1000), 0, SimplexDt{n: 0.0, a: 0.05, freq: 0.00005, oct: 8}))
		simplexNoise6 := float32(Noise3dSimplex(float64(pt6[0]+1000), float64(pt6[1]+25000), float64(pt6[2]+1000), 0, SimplexDt{n: 0.0, a: 0.05, freq: 0.00005, oct: 8}))

		pt1 = mgl32.Vec3{pt1[0] * (1 + simplexNoise1), pt1[1] * (1 + simplexNoise1), pt1[2] * (1 + simplexNoise1)}
		pt2 = mgl32.Vec3{pt2[0] * (1 + simplexNoise2), pt2[1] * (1 + simplexNoise2), pt2[2] * (1 + simplexNoise2)}
		pt3 = mgl32.Vec3{pt3[0] * (1 + simplexNoise3), pt3[1] * (1 + simplexNoise3), pt3[2] * (1 + simplexNoise3)}
		pt4 = mgl32.Vec3{pt4[0] * (1 + simplexNoise4), pt4[1] * (1 + simplexNoise4), pt4[2] * (1 + simplexNoise4)}
		pt5 = mgl32.Vec3{pt5[0] * (1 + simplexNoise5), pt5[1] * (1 + simplexNoise5), pt5[2] * (1 + simplexNoise5)}
		pt6 = mgl32.Vec3{pt6[0] * (1 + simplexNoise6), pt6[1] * (1 + simplexNoise6), pt6[2] * (1 + simplexNoise6)}

		rgb := (simplexNoise1 + simplexNoise4 + simplexNoise6) / 3
		var water float32 = 0.0
		var ground float32 = 0.0
		if rgb < 0.3 {
			water = 0.75
		} else {
			ground = 0.25
		}
		newPoints = append(newPoints, pt1[0], pt1[1], pt1[2], simplexNoise1, simplexNoise1+ground, simplexNoise1+water)
		newPoints = append(newPoints, pt4[0], pt4[1], pt4[2], simplexNoise4, simplexNoise4+ground, simplexNoise4+water)
		newPoints = append(newPoints, pt6[0], pt6[1], pt6[2], simplexNoise6, simplexNoise6+ground, simplexNoise6+water)

		water = 0
		ground = 0
		rgb = float32(simplexNoise4+simplexNoise5+simplexNoise6) / 3
		if rgb < 0.3 {
			water = 0.75
		} else {
			ground = 0.25
		}
		newPoints = append(newPoints, pt4[0], pt4[1], pt4[2], simplexNoise4, simplexNoise4+ground, simplexNoise4+water)
		newPoints = append(newPoints, pt5[0], pt5[1], pt5[2], simplexNoise5, simplexNoise5+ground, simplexNoise5+water)
		newPoints = append(newPoints, pt6[0], pt6[1], pt6[2], simplexNoise6, simplexNoise6+ground, simplexNoise6+water)

		water = 0
		ground = 0
		rgb = float32(simplexNoise5+simplexNoise6+simplexNoise3) / 3
		if rgb < 0.3 {
			water = 0.75
		} else {
			ground = 0.25
		}
		newPoints = append(newPoints, pt5[0], pt5[1], pt5[2], simplexNoise5, simplexNoise5+ground, simplexNoise5+water)
		newPoints = append(newPoints, pt6[0], pt6[1], pt6[2], simplexNoise6, simplexNoise6+ground, simplexNoise6+water)
		newPoints = append(newPoints, pt3[0], pt3[1], pt3[2], simplexNoise3, simplexNoise3+ground, simplexNoise3+water)

		water = 0
		ground = 0
		rgb = float32(simplexNoise4+simplexNoise2+simplexNoise5) / 3
		if rgb < 0.3 {
			water = 0.75
		} else {
			ground = 0.25
		}
		newPoints = append(newPoints, pt4[0], pt4[1], pt4[2], simplexNoise4, simplexNoise4+ground, simplexNoise4+water)
		newPoints = append(newPoints, pt2[0], pt2[1], pt2[2], simplexNoise2, simplexNoise2+ground, simplexNoise2+water)
		newPoints = append(newPoints, pt5[0], pt5[1], pt5[2], simplexNoise5, simplexNoise5+ground, simplexNoise5+water)
	}
	return newPoints
}

func main() {
	rand.Seed(time.Now().UnixNano())
	NoiseInitPermtables(42)

	/*
		planet1 := PlanetData{
			seaLevel:    0.325,
			totalPoints: 10000,
			scale:       500,
			dt:          SimplexDt{n: 0.0, a: 0.25, freq: 0.75, oct: 5},
			dtBillow:    SimplexDt{n: 0.0, a: 0.25, freq: 0.75, oct: 1},
			colorGround: ColorRGB{40, 100, 23},
			colorSea:    ColorRGB{0, 62, 120},
			pos:         mgl32.Mat3{-2000, 0, 0},
		}
		planet2 := PlanetData{
			seaLevel:    0.325,
			totalPoints: 10000,
			scale:       750,
			dt:          SimplexDt{n: 0.0, a: 0.75, freq: 0.15, oct: 8},
			dtBillow:    SimplexDt{n: 0.0, a: 0.45, freq: 0.75, oct: 1},
			colorGround: ColorRGB{25, 25, 25},
			colorSea:    ColorRGB{0, 100, 100},
			pos:         mgl32.Mat3{-6000, 0, 0},
		}
		planet3 := PlanetData{
			seaLevel:    0.325,
			totalPoints: 5000,
			scale:       250,
			dt:          SimplexDt{n: 0.0, a: 1.0, freq: 2.5, oct: 1},
			dtBillow:    SimplexDt{n: 0.0, a: 0.45, freq: 2.5, oct: 1},
			colorGround: ColorRGB{100, 0, 0},
			colorSea:    ColorRGB{200, 200, 200},
			pos:         mgl32.Mat3{-4000, -2500, 0},
		}

		generatePlanet(&points, planet1)
		generatePlanet(&points, planet2)
		generatePlanet(&points, planet3)
	*/
	points = []float32{
		-7500, 0, 0, 1, 0, 0,
		0, 7500, 0, 0, 1, 0,
		0, 0, 7500, 0, 0, 1,

		0, 0, 7500, 0, 0, 1,
		0, 7500, 0, 0, 1, 0,
		7500, 0, 0, 1, 0, 1,

		7500, 0, 0, 1, 0, 1,
		0, 7500, 0, 0, 1, 0,
		0, 0, -7500, 1, 1, 0,

		0, 0, -7500, 1, 1, 0,
		0, 7500, 0, 0, 1, 0,
		-7500, 0, 0, 1, 0, 0,

		-7500, 0, 0, 1, 0, 0,
		0, -7500, 0, 0, 1, 1,
		0, 0, 7500, 0, 0, 1,

		0, 0, 7500, 0, 0, 1,
		0, -7500, 0, 0, 1, 1,
		7500, 0, 0, 1, 0, 1,

		7500, 0, 0, 1, 0, 1,
		0, -7500, 0, 0, 1, 1,
		0, 0, -7500, 1, 1, 0,

		0, 0, -7500, 1, 1, 0,
		0, -7500, 0, 0, 1, 1,
		-7500, 0, 0, 1, 0, 0,
	}

	//points = append(points, (0), (0), (0), 1, 0, 0)

	var cbo Cbo
	var colorRColorRGB ColorRGB
	colorRColorRGB.r, colorRColorRGB.g, colorRColorRGB.b = 1.0, 1.0, 1.0
	cbo.pos[0], cbo.pos[1], cbo.pos[2] = 0.0, 0.0, -10000.0

	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	redUniform := gl.GetUniformLocation(program, gl.Str("red\x00"))

	cameraId := gl.GetUniformLocation(program, gl.Str("camera\x00"))

	gl.PointSize(5)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	vao := makeVao(points)

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
		EventsKeyboard(&cbo, &colorRColorRGB, &keys, &gameValues)
		EventsMouse(&cbo)
		setCamera(cameraId, &cbo)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		if keys.v == "active" {
			points = subdivideTriangles(points)
			vao = makeVao(points)
			fmt.Println(len(points) / 6 / 3)
		}

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points)/3))

		glfw.PollEvents()
		window.SwapBuffers()

		gl.Uniform1f(redUniform, colorRColorRGB.r)
	}
}

func setCamera(cameraId int32, cbo *Cbo) {
	camera := mgl32.HomogRotate3D(cbo.rot.X(), mgl32.Vec3{1, 0, 0})
	camera = camera.Mul4(mgl32.HomogRotate3D(cbo.rot.Y(), mgl32.Vec3{0, 1, 0}))
	camera = camera.Mul4(mgl32.Translate3D(cbo.pos.X(), cbo.pos.Y(), cbo.pos.Z()))

	projection := mgl32.Perspective(mgl32.DegToRad(80.0), float32(width)/float32(height), 0.1, 50000)
	view := projection.Mul4(camera)

	gl.UniformMatrix4fv(cameraId, 1, false, &view[0])
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 6*4, 0)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 6*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return vao
}

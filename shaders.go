package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func _getShaderSource(file string) string {
	data, error := ioutil.ReadFile(file)
	if error != nil {
		fmt.Println("Could not open", file)
		os.Exit(1)
	}
	return string(data) + "\x00"
}

func GetCameraMatrix(cbo Cbo) mgl32.Mat4 {
	camera := mgl32.HomogRotate3D(cbo.rot.X(), mgl32.Vec3{1, 0, 0})
	camera = camera.Mul4(mgl32.HomogRotate3D(cbo.rot.Y(), mgl32.Vec3{0, 1, 0}))
	camera = camera.Mul4(mgl32.Translate3D(cbo.pos.X(), cbo.pos.Y(), cbo.pos.Z()))

	return camera
}

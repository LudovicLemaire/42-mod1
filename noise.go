package main

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Vector2 struct {
	x float64
	y float64
}

type SimplexDt struct {
	n    float64
	a    float64
	freq float64
	oct  int
}

type Seeder struct {
	perm  [512]int
	gradP [512]Vector3
}

var grad3 = [12]Vector3{
	{1, 1, 0}, {-1, 1, 0}, {1, -1, 0}, {-1, -1, 0},
	{1, 0, 1}, {-1, 0, 1}, {1, 0, -1}, {-1, 0, -1},
	{0, 1, 1}, {0, -1, 1}, {0, 1, -1}, {0, -1, -1},
}

func (vec3 Vector3) dot3(other Vector3) float64 {
	return vec3.x*other.x + vec3.y*other.y + vec3.z*other.z
}

func (vec3 Vector3) dot2(other Vector2) float64 {
	return vec3.x*other.x + vec3.y*other.y
}

var Seeds [5]Seeder

var defaultPermutations = [512]int{151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180}

func Noise3dSimplex(x, y, z float64, seed int, dt SimplexDt) float64 {
	x = math.Abs(x)
	y = math.Abs(y)
	z = math.Abs(z)

	// var n float64 = 0.0
	// var a float64 = 0.25
	// var freq float64 = 0.75

	for currOctave := 0; currOctave < dt.oct; currOctave++ {
		var v float64 = dt.a * simplex3d(float64(x)*dt.freq, float64(y)*dt.freq, float64(z)*dt.freq, seed)
		dt.n += v

		dt.a *= 0.5
		dt.freq *= 2.0
	}

	dt.n = (dt.n + 1.0) * 0.5

	if dt.n < 0 {
		dt.n = 0
	}
	return dt.n
}

func simplex3d(xin, yin, zin float64, seed int) float64 {
	var F3 float64 = 1.0 * 0.333333333
	var G3 float64 = 1.0 * 0.166666667
	var n0, n1, n2, n3 float64 // Noise contributions from the four corners

	// Skew the input space to determine which simplex cell we're in
	var s float64 = (xin + yin + zin) * F3 // Hairy factor for 2D
	var i = int(xin + s)
	var j = int(yin + s)
	var k = int(zin + s)

	var t = float64(i+j+k) * G3
	var x0 = xin - float64(i) + t // The x,y distances from the cell origin, unskewed.
	var y0 = yin - float64(j) + t
	var z0 = zin - float64(k) + t

	// For the 3D case, the simplex shape is a slightly irregular tetrahedron.
	// Determine which simplex we are in.
	var i1, j1, k1 int // Offsets for second corner of simplex in (i,j,k) coords
	var i2, j2, k2 int // Offsets for third corner of simplex in (i,j,k) coords
	if x0 >= y0 {
		if y0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0
		} else if x0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 0
			k2 = 1
		} else {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 1
			j2 = 0
			k2 = 1
		}
	} else {
		if y0 < z0 {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 0
			j2 = 1
			k2 = 1
		} else if x0 < z0 {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 0
			j2 = 1
			k2 = 1
		} else {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0
		}
	}
	// A step of (1,0,0) in (i,j,k) means a step of (1-c,-c,-c) in (x,y,z),
	// a step of (0,1,0) in (i,j,k) means a step of (-c,1-c,-c) in (x,y,z), and
	// a step of (0,0,1) in (i,j,k) means a step of (-c,-c,1-c) in (x,y,z), where
	// c = 1/6.
	var x1 = x0 - float64(i1) + G3 // Offsets for second corner
	var y1 = y0 - float64(j1) + G3
	var z1 = z0 - float64(k1) + G3

	var x2 = x0 - float64(i2) + 2*G3 // Offsets for third corner
	var y2 = y0 - float64(j2) + 2*G3
	var z2 = z0 - float64(k2) + 2*G3

	var x3 = x0 - 1 + 3*G3 // Offsets for fourth corner
	var y3 = y0 - 1 + 3*G3
	var z3 = z0 - 1 + 3*G3

	// Work out the hashed gradient indices of the four simplex corners
	i &= 255
	j &= 255
	k &= 255

	// Calculate the contribution from the four corners
	var t0 = 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * Seeds[seed].gradP[i+Seeds[seed].perm[j+Seeds[seed].perm[k]]].dot3(Vector3{x0, y0, z0}) // (x,y) of grad3 used for 2D gradient
	}
	var t1 = 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * Seeds[seed].gradP[i+i1+Seeds[seed].perm[j+j1+Seeds[seed].perm[k+k1]]].dot3(Vector3{x1, y1, z1})
	}
	var t2 = 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * Seeds[seed].gradP[i+i2+Seeds[seed].perm[j+j2+Seeds[seed].perm[k+k2]]].dot3(Vector3{x2, y2, z2})
	}
	var t3 = 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0 {
		n3 = 0
	} else {
		t3 *= t3
		n3 = t3 * t3 * Seeds[seed].gradP[i+1+Seeds[seed].perm[j+1+Seeds[seed].perm[k+1]]].dot3(Vector3{x3, y3, z3})
	}
	// Add contributions from each corner to get the final noise value.
	// The result is scaled to return values in the interval [-1,1].
	return 32 * (n0 + n1 + n2 + n3)
}

func _seedFunc(seed int) Seeder {
	var perm [512]int
	var gradP [512]Vector3

	if seed > 0 && seed < 1 {
		// Scale the seed out
		seed *= 65536
	}

	seed = int(seed)
	if seed < 256 {
		seed |= seed << 8
	}

	for i := 0; i < 256; i++ {
		var v int
		if i&1 == 1 {
			v = defaultPermutations[i] ^ (seed & 255)
		} else {
			v = defaultPermutations[i] ^ ((seed >> 8) & 255)
		}

		perm[i] = v
		perm[i+256] = v
		gradP[i] = grad3[v%12]
		gradP[i+256] = grad3[v%12]
	}
	return Seeder{perm, gradP}
}

func NoiseInitPermtables(seed float64) {
	// Generate table permutations with seed \\
	Seeds[0] = _seedFunc(int(seed * 1.2))
	Seeds[1] = _seedFunc(int(seed * 2.5))
	Seeds[2] = _seedFunc(int(seed * 3.6))
	Seeds[3] = _seedFunc(int(seed * 4.7))
	Seeds[4] = _seedFunc(int(seed * 2))
}

func simplex2d(xin, yin float64, seed int) float64 {
	// Skewing and unskewing factors for 2, 3, and 4 dimensions
	var F2 = 0.5 * (math.Sqrt(3) - 1)
	var G2 = (3 - math.Sqrt(3)) * 0.166666667

	var n0, n1, n2 float64 // Noise contributions from the three corners
	// Skew the input space to determine which simplex cell we are in
	var s float64 = (xin + yin) * F2 // Hairy factor for 2D
	var i int = int(xin + s)
	var j int = int(yin + s)
	var t float64 = float64(i+j) * G2
	var x0 float64 = xin - float64(i) + t // The x, y distances from the cell origin
	var y0 float64 = yin - float64(j) + t
	// Determine which simplex we are in
	var i1, j1 float64 // Offsets for second (middle) corner of simplex in (i, j) coords
	if x0 > y0 {       // lower triangle, XY order: (0, 0) -> (1, 0) -> (1, 1)
		i1 = 1
		j1 = 0
	} else { // upper triangle, YX order: (0, 0) -> (0, 1) -> (1, 1)
		i1 = 0
		j1 = 1
	}
	// A step of (1,0) in (i,j) means a step of (1-c,-c) in (x,y), and
	// a step of (0,1) in (i,j) means a step of (-c,1-c) in (x,y), where
	// c = (3-sqrt(3))/6
	var x1 = x0 - i1 + G2 // Offsets for middle corner in (x,y)
	var y1 = y0 - j1 + G2
	var x2 = x0 - 1 + 2*G2 // Offsets for last corner in (x,y)
	var y2 = y0 - 1 + 2*G2
	// Work out the hashed gradient indices of the three simplex corners
	i &= 255
	j &= 255

	// Calculate the contribution from the three corners
	var t0 = 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * Seeds[seed].gradP[i+Seeds[seed].perm[j]].dot2(Vector2{x0, y0}) // (x,y) of grad3 used for 2D gradient
	}
	var t1 = 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * Seeds[seed].gradP[i+int(i1)+Seeds[seed].perm[j+int(j1)]].dot2(Vector2{x1, y1})
	}
	var t2 = 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * Seeds[seed].gradP[i+1+Seeds[seed].perm[j+1]].dot2(Vector2{x2, y2})
	}
	// Add contributions from each corner to get the final noise value
	// The result is scaled to return values in the interval [-1, 1]
	return 70 * (n0 + n1 + n2)
}

func Noise2dSimplex(x, y, n, a, freq float64, octave, seed int) float64 {
	x = math.Abs(x)
	y = math.Abs(y)

	for octave := 0; octave < 8; octave++ {
		var v float64 = a * simplex2d(float64(x)*freq, float64(y)*freq, seed)
		n += v

		a *= 0.5
		freq *= 2.0
	}

	n = (n + 1.0) * 0.5

	if n < 0 {
		n = 0
	} else if n >= 1 {
		n = 1
	}

	return n
}

func generateVarianteColorBillowNoise(rgb mgl32.Vec3, x, y, z float64) mgl32.Vec3 {
	colorVariationPerlin := float32(Noise3dSimplex(x, y, z, 0.0, SimplexDt{0.0, 1.0, 0.035, 5})) / 2
	var newColorVariationR float32
	var newColorVariationG float32
	var newColorVariationB float32
	if rgb[0]-colorVariationPerlin < 0 {
		newColorVariationR = 0
	} else {
		newColorVariationR = (rgb[0]) - (colorVariationPerlin)
	}
	if rgb[1]-colorVariationPerlin < 0 {
		newColorVariationG = 0
	} else {
		newColorVariationG = (rgb[1]) - (colorVariationPerlin)
	}
	if rgb[2]-colorVariationPerlin < 0 {
		newColorVariationB = 0
	} else {
		newColorVariationB = (rgb[2]) - (colorVariationPerlin)
	}
	if math.IsNaN(float64(newColorVariationR)) {
		fmt.Println(colorVariationPerlin)
		fmt.Println(x, y, z)
		fmt.Println(mgl32.Vec3{newColorVariationR, newColorVariationG, newColorVariationB})
	}
	return mgl32.Vec3{newColorVariationR, newColorVariationG, newColorVariationB}
}

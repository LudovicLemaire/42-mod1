#version 400
layout(location = 0) in vec3 vertex_position;
layout(location = 1) in vec3 vertex_color;

uniform mat4 camera;
uniform float white;
out vec4 theColor;

void main() {
	gl_Position = camera * vec4(vertex_position.x, vertex_position.y, vertex_position.z, 1.0);
	if (white == 0.0) {
		theColor = vec4(vertex_color, 1.0);
	} else {
		theColor = vec4(1.0, 1.0, 1.0, 1.0);
	}
	
}
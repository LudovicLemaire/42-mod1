# Golang-OpenGL-template
Simple Golang/OpenGL template I did so I can reuse it when I need to test so graphical representation, mostly for school projects.

## Excecute
```
go run .
```
## Requirement
* Golang
* OpenGL 4.1
- - -
## It has
* a complete event handling, with different states (click, hold, released, null).
* a camera to move around
* a buffer ready to draw triangle/color from a slice.

![subdivide](.git_resources/octahedron.gif)

* buffer modification in real-time, here an example by subdividing triangles into 4. Where noise is applyied to render a planet.

![subdivide](.git_resources/subdivide2.gif)
![subdivide](.git_resources/subdivide.gif)

`/!\
Overall code isn't cleaned nor commented.
It's just a quick and dirty way for me to test stuff.
/!\`
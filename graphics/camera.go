package graphics;

import "gotracer/vmath";

// Camera object describes how the objects are projected into the screen
type Camera struct {
	// Position of the camera in the world
	Origin *vmath.Vector3;

	//The Lower left corner of the camera relative to the center considering the vertical and horizontal sizes.
	LowerLeftCorner *vmath.Vector3;

	// Vertical size of the camera (usually only uses Y)
	Vertical *vmath.Vector3;

	// Horizontal size of the camera (usually only uses X)
	Horizontal *vmath.Vector3;
}

// Create a new camera with default values.
func NewCamera() *Camera {
	var c = new(Camera);
	c.LowerLeftCorner = vmath.NewVector3(-2.0, -1.0, -1.0);
	c.Horizontal = vmath.NewVector3(4.0, 0.0, 0.0);
	c.Vertical = vmath.NewVector3(0.0, 2.0, 0.0);
	c.Origin = vmath.NewVector3(0.0, 0.0, 0.0);
	return c;
}

// Get a ray from this camera, from a normalized UV screen coordinate.
func (s *Camera) GetRay(u float64, v float64) bool {

	//TODO <ADD CODE HERE>

	return false;
}


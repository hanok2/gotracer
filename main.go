package main;

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"gotracer/graphics"
	"gotracer/hitable"
	"gotracer/vmath"
);

var origin = vmath.NewVector3(0.0, 0.0, 0.0);
var world hitable.HitableList;
var camera *graphics.Camera;

func run() {
	var width = 640;
	var height = 480;

	var bounds = pixel.R(0, 0, float64(width), float64(height));

	var config = pixelgl.WindowConfig{
		Resizable: false,
		Undecorated: false,
		VSync: false,
		Title: "Gotracer",
		Bounds: bounds};

	var window, err = pixelgl.NewWindow(config);
	
	CheckError(err);

	for !window.Closed() {
		
		var start = time.Now();

		window.Clear(colornames.Black);

		var picture = Raytrace(bounds);

		var sprite = pixel.NewSprite(picture, picture.Bounds());
		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()));	

		window.Update();

		var delta = time.Since(start);
		log.Printf("Frame time %s", delta);
	}
}

func main() {
	camera = graphics.NewCamera();

	// Prepare the scene
	world.Add(hitable.NewSphere(0.5, vmath.NewVector3(0.0, 0.0, -1.0)));
	world.Add(hitable.NewSphere(100.0, vmath.NewVector3(0.0, -100.5, -1.0)));

	// Start the renderer
	pixelgl.Run(run)
}

func CalculateColor(r *vmath.Ray) *vmath.Vector3 {
	var rec = hitable.NewHitRecord();

	if world.Hit(r, 0.0, 9999999.0, rec) {
		// Normal color
		var color = vmath.NewVector3(rec.Normal.X + 1.0, rec.Normal.Y + 1.0, rec.Normal.Z + 1.0);
		color.MulScalar(0.5);
		return color;

	} else {
		// Background
		var unitDirection = r.Direction.UnitVector();
		var t = 0.5 * (unitDirection.Y + 1.0);

		var a = vmath.NewVector3(1.0, 1.0, 1.0);
		a.MulScalar(1.0 - t);

		var b = vmath.NewVector3(0.5, 0.7, 1.0);
		b.MulScalar(t);

		a.Add(b);

		return a;
	}
}

//Render sky with raytrace
func Raytrace(bounds pixel.Rect) *pixel.PictureData {

	var size = bounds.Size();
	var aspect = size.X / size.Y;
	var scale = 2.0;

	var lowerLeftCorner = vmath.NewVector3(-scale / 2.0 * aspect, -1.0, -1.0);
	var vertical = vmath.NewVector3(0.0, scale, 0.0);
	var horizontal = vmath.NewVector3(scale * aspect, 0.0, 0.0);
	var picture = pixel.MakePictureData(bounds);

	var nx int = int(size.X);
	var ny int = int(size.Y);

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			var u = float64(i) / size.X;
			var v = float64(j) / size.Y;

			var hor = horizontal.Clone();
			hor.MulScalar(u);

			var vert = vertical.Clone();
			vert.MulScalar(v);

			var direction = lowerLeftCorner.Clone();
			direction.Add(hor);
			direction.Add(vert);

			var ray = vmath.NewRay(origin, direction);

			//Calculate color
			var color = CalculateColor(ray);

			var ir int = int(255 * color.X);
			var ig int = int(255 * color.Y);
			var ib int = int(255 * color.Z);

			//Write to picture
			var index = picture.Index(pixel.Vec{X:float64(i), Y:float64(j)});
			picture.Pix[index].R = uint8(ir);
			picture.Pix[index].G = uint8(ig);
			picture.Pix[index].B = uint8(ib);
		}
	}

	return picture;
}

//Dot product between two vectors
func Dot(a *vmath.Vector3, b *vmath.Vector3) float64 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z;
}

//Cross product between two vectors
func Cross(result *vmath.Vector3, a *vmath.Vector3, b *vmath.Vector3) {
	result.X = a.Y * b.Z - a.Z * b.Y;
	result.Y = a.Z * b.X - a.X * b.Z;
	result.Z = a.X * b.Y - a.Y * b.X;
}


// Write the frame to a PPM file string
func WritePPM(picture *pixel.PictureData, fname string) {
	var size = picture.Rect.Size();

	var nx int = int(size.X);
	var ny int = int(size.Y);
	
	var file, err = os.Create("sky.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			//Write to file
			var index = picture.Index(pixel.Vec{X:float64(i), Y:float64(j)});
			file.WriteString(strconv.Itoa(int(picture.Pix[index].R)) + " " + strconv.Itoa(int(picture.Pix[index].G)) + " " + strconv.Itoa(int(picture.Pix[index].B)) + "\n");
		}
	}

	//Close file
	file.Sync();
	file.Close();
}

//CheckError an error	
func CheckError(e error) {
	if e != nil {
		panic(e);
	}
}

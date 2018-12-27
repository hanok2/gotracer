package main;

import "os";
import "strconv";
//import "fmt";
import "math";
//import "bufio"
//import "io/ioutil"

import "github.com/faiface/pixel";
import "github.com/faiface/pixel/pixelgl";
import "golang.org/x/image/colornames";

import "gotracer/vmath";
//import "gotracer/graphics";
//import "gotracer/hitable";

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

	window.Clear(colornames.Black);

	var picture = pixel.MakePictureData(bounds);

	var a uint8 = 0;

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var index = picture.Index(pixel.Vec{X:float64(x), Y:float64(y)});
			picture.Pix[index].R = a;
			picture.Pix[index].G = a;
			picture.Pix[index].B = a;
		}
		
		a++;
		if a > 255 {
			a = 0;
		}

	}

	var sprite = pixel.NewSprite(picture, picture.Bounds());
	sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()));

	//First raytrace
	RenderSky();

	for !window.Closed() {
		window.Update()
	}
}

func main() {
	//Create window


	//var bounds = new(pixelgl);
	//config.Bounds = ;

	pixelgl.Run(run)
}

func GetColor(r *vmath.Ray) *vmath.Vector3 {
	var t = HitSphere(vmath.NewVector3(0.0, 0.0, -1.0), 0.5, r);

	if(t > 0.0) {
		var n = r.PointAtParameter(t);
		n.Sub(vmath.NewVector3(0, 0, -1));
		n.Normalize();
		n.Add(vmath.NewVector3(1.0, 1.0, 1.0));
		n.MulScalar(0.5);
		return n;
	}

	var unitDirection = r.Direction.UnitVector();
	t = 0.5 * (unitDirection.Y + 1.0);

	var a = vmath.NewVector3(1.0, 1.0, 1.0);
	a.MulScalar(1.0 - t);

	var b = vmath.NewVector3(0.5, 0.7, 1.0);
	b.MulScalar(t);

	a.Add(b);
	
	return a;
}

func HitSphere(center *vmath.Vector3, radius float64, ray *vmath.Ray) float64 {
	var oc = ray.Origin.Clone();
	oc.Sub(center);

	var a = Dot(ray.Direction, ray.Direction);
	var b = 2.0 * Dot(oc, ray.Direction);
	var c = Dot(oc, oc) - radius * radius;
	var discriminant = b * b - 4 * a * c;

	if(discriminant < 0) {
		return -1.0;
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a);
	}
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

//Render sky with raytrace
func RenderSky() {
	var nx int = 200;
	var ny int = 100;

	var lowerLeftCorner = vmath.NewVector3(-2.0, -1.0, -1.0);
	var horizontal = vmath.NewVector3(4.0, 0.0, 0.0);
	var vertical = vmath.NewVector3(0.0, 2.0, 0.0);
	var origin = vmath.NewVector3(0.0, 0.0, 0.0);

	var file, err = os.Create("sky.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			var u = float64(i) / float64(nx);
			var v = float64(j) / float64(ny);

			var hor = horizontal.Clone();
			hor.MulScalar(u);

			var vert = vertical.Clone();
			vert.MulScalar(v);

			var direction = lowerLeftCorner.Clone();
			direction.Add(hor);
			direction.Add(vert);

			var ray = vmath.NewRay(origin, direction);

			//Calculate color
			var color = GetColor(ray);

			var ir int = int(255 * color.X);
			var ig int = int(255 * color.Y);
			var ib int = int(255 * color.Z);

			file.WriteString(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
		}
	}

	//Close file
	file.Sync();
	file.Close();
}

//RenderGradient the image
func RenderGradient() {
	var nx int = 200;
	var ny int = 100;

	var file, err = os.Create("gradient.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			//Calculate color
			var color = vmath.NewVector3(float64(i) / float64(nx), float64(j) / float64(ny), 0.2);

			var ir int = int(256 * color.X);
			var ig int = int(256 * color.Y);
			var ib int = int(256 * color.Z);

			file.WriteString(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
		}
	}

	//Close file
	file.Sync();
	file.Close();
}

// Write the frame to a PPM file string
func WritePPM(picture *pixel.PictureData, fname string) string {
	
	/*var file, err = os.Create("sky.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			var ir int = int(255 * color.X);
			var ig int = int(255 * color.Y);
			var ib int = int(255 * color.Z);

			file.WriteString(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
		}
	}

	//Close file
	file.Sync();
	file.Close();*/

	return "TODO";
}

//CheckError an error	
func CheckError(e error) {
	if e != nil {
		panic(e);
	}
}

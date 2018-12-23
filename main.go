package main;

import "fmt";
import "os";
import "strconv";
//import "bufio"
//import "io/ioutil"

//import "github.com/faiface/pixel/pixelgl";

func main() {
	Render();
}

//Render the image
func Render() {
	var nx int = 640;
	var ny int = 480;

	var file, err = os.Create("output.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");
	fmt.Println("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			//Calculate color
			var r float64 = float64(i) / float64(nx);
			var g float64 = float64(j) / float64(ny);
			var b float64 = 0.2;

			var ir int = int(256 * r);
			var ig int = int(256 * g);
			var ib int = int(256 * b);

			file.WriteString(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
			fmt.Println(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
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
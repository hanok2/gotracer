package hitable;

import "gotracer/vmath";
import "math"

// Sphere is hitable object represented by a center point and a radius.
type Sphere struct {
	Radius float64;
	Center *vmath.Vector3;
}

func (s *Sphere) Hit(ray *vmath.Ray, tmin float64, tmax float64, rec *HitRecord) bool {

	var oc = ray.Origin.Clone();
	oc.Sub(s.Center);

	var a = Dot(ray.Direction, ray.Direction);
	var b = 2.0 * Dot(oc, ray.Direction);
	var c = Dot(oc, oc) - s.Radius * s.Radius;
	var discriminant = b * b - a * c;

	if(discriminant > 0) {

		//First root result
		var temp float64 = (-b - math.Sqrt(discriminant)) / a;

		if(temp < tmax && temp > tmin) {
			rec.T = temp;
			rec.P = ray.PointAtParameter(temp);
			rec.Normal = rec.P.Clone();
			rec.Normal.Sub(s.Center);
			rec.Normal.DivideScalar(s.Radius);
			return true;
		}

		//Second root result
		temp = (-b + math.Sqrt(discriminant)) / a;
		
		if(temp < tmax && temp > tmin) {
			rec.T = temp;
			rec.P = ray.PointAtParameter(temp);
			rec.Normal = rec.P.Clone();
			rec.Normal.Sub(s.Center);
			rec.Normal.DivideScalar(s.Radius);
			return true;
		}
		
	}

	return false;
}


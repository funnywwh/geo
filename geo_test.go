package geo

import (
	"math"
	"testing"
)

func TestPointEquals(t *testing.T) {
	pt1 := &Point{
		Lat: 22.3,
		Lng: 22.4,
	}
	pt2 := &Point{
		Lat: 22.3,
		Lng: 22.4,
	}
	if pt1.Equals(pt2) == false {
		t.Errorf("pt1:%+v != pt2:%+v", pt1, pt2)
		t.FailNow()
	}
	pt3 := &Point{
		Lat: 22.3,
		Lng: 22.6,
	}
	if pt1.Equals(pt3) == true {
		t.Errorf("pt1:%+v == pt3:%+v", pt1, pt3)
		t.FailNow()
	}
}

func TestGetBounds(t *testing.T) {
	ply := &Polyline{
		Points: []*Point{
			&Point{
				Lat: 22.3,
				Lng: 22.4,
			},
			&Point{
				Lat: 22,
				Lng: 42.4,
			},
			&Point{
				Lat: 12.3,
				Lng: 32.4,
			},
			&Point{
				Lat: 22.3,
				Lng: 82.4,
			},
		},
	}
	rt := ply.GetBounds()
	if rt.NorthEast.Equals(&Point{
		Lat: 22.3,
		Lng: 82.4,
	}) == false || rt.SouthWest.Equals(&Point{
		Lat: 12.3,
		Lng: 22.4,
	}) == false {
		t.Errorf("rt:%+v", rt)
		t.FailNow()
	}
}

func TestGetDistance(t *testing.T) {
	pt1 := &Point{
		Lat: 22,
		Lng: 82,
	}
	pt2 := &Point{
		Lat: 22,
		Lng: 82,
	}
	dis := GetDistance(pt1, pt2)
	if dis > 0 {
		t.Errorf("dis:%v\n", dis)
	}
	pt3 := &Point{
		Lat: 102,
		Lng: 82,
	}
	dis = GetDistance(pt3, pt1)
	if dis < 5.782133290364891e+06 {
		t.Error(dis)
	}
}

func TestIsPointInRect(t *testing.T) {
	pt1 := &Point{
		Lat: 22,
		Lng: 82,
	}
	rt := &Rect{
		SouthWest: Point{
			Lat: 22,
			Lng: 84,
		},
		NorthEast: Point{
			Lat: 25,
			Lng: 82,
		},
	}
	in := IsPointInRect(pt1, rt)
	if in {
		t.Error("pt:%v in rt:%v", pt1, rt)
	}
	rt = &Rect{
		SouthWest: Point{
			Lat: 22,
			Lng: 80,
		},
		NorthEast: Point{
			Lat: 25,
			Lng: 84,
		},
	}
	in = IsPointInRect(pt1, rt)
	if in == false {
		t.Error(in)
	}
}
func TestIsPointInCircle(t *testing.T) {
	var circle Circle
	circle.Center = Point{Lat: 22, Lng: 80}
	circle.Radius = 100
	if IsPointInCircle(&circle.Center, &circle) == false {
		t.FailNow()
	}
	var pt2 = Point{Lat: 22.001, Lng: 80.001}
	d := GetDistance(&pt2, &circle.Center)
	t.Logf("d2:%f", d)
	var pt3 = Point{Lat: 22.0001, Lng: 80.0007}
	d = GetDistance(&pt3, &circle.Center)
	t.Logf("d3:%f", d)
	if IsPointInCircle(&pt2, &circle) == true {
		t.FailNow()
		t.Fatalf("pt2 fail")
	}
	if IsPointInCircle(&pt3, &circle) == false {
		t.FailNow()
		t.Fatalf("pt3 fail")
	}
}

func TestIsPointOnPolyline(t *testing.T) {
	p1 := Point{Lat: 40.059295, Lng: 116.300306}
	p2 := Point{Lat: 40.059019, Lng: 116.310906}
	p3 := Point{Lat: 40.054242, Lng: 116.303073}
	pl := Polyline{
		Points: []*Point{&p1, &p2, &p3},
	}
	outpt := Point{Lat: 40.061028, Lng: 116.303154}
	onpt := Point{Lat: 40.059137, Lng: 116.306316}
	if IsPointInPolygon(&outpt, &pl) == true {
		t.Fatalf("fail:outpt:%v should not in pl:%v", outpt, pl)
	}
	if IsPointInPolygon(&onpt, &pl) == false {
		t.Fatalf("fail onpt:%v should on pl:%v", onpt, pl)
	}
}
func TestGetPolylineDistance(t *testing.T) {
	p1 := Point{Lat: 40.059295, Lng: 116.300306}
	p2 := Point{Lat: 40.059019, Lng: 116.310906}
	p3 := Point{Lat: 40.054242, Lng: 116.303073}
	d := GetDistance(&p1, &p2)
	d += GetDistance(&p2, &p3)
	pl := Polyline{
		Points: []*Point{&p1, &p2, &p3},
	}
	if math.Abs(d-GetPolylineDistance(&pl)) > 0.0000001 {
		t.Fatalf("d-GetPolylineDistance(p1) > 0.0001")
	}
}

func TestGetPolygonArea(t *testing.T) {

	var pt1 = Point{116.395, 39.910}
	var pt2 = Point{116.394, 39.918}
	var pt3 = Point{116.396, 39.919}
	var pt4 = Point{116.404, 39.920}
	var pt5 = Point{116.406, 39.913}
	pl := Polyline{
		Points: []*Point{&pt1, &pt2, &pt3, &pt4, &pt5},
	}
	area := GetPolygonArea(&pl)
	t.Logf("Polyline area:%f", area)

	if math.Abs(area-810876.60) > 0.1 {
		t.Fatalf("area - 810876.60 > 0.01")
	}
}

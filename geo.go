package geo

import (
	"math"
)

const (
	EARTHRADIUS = 6370996.81
	F_ZERO      = 2e-10
)

type Point struct {
	Lng float64 //经度
	Lat float64 //纬度
}

func (this *Point) Equals(p2 *Point) bool {
	if math.Abs(this.Lat-p2.Lat) <= F_ZERO &&
		math.Abs(this.Lng-p2.Lng) <= F_ZERO {
		return true
	}
	return false
}

type Rect struct {
	SouthWest Point
	NorthEast Point
}

type Circle struct {
	Center Point
	Radius float64
}

type Polyline struct {
	Points []*Point
}

/**
 * 计算多边区的外包含矩形
 */
func (this *Polyline) GetBounds() (rt *Rect) {
	var min, max float64
	min = math.MaxFloat64
	max = -math.MaxFloat64
	for _, v := range this.Points {
		min = math.Min(min, v.Lat)
		max = math.Max(max, v.Lat)
	}
	rt.SouthWest.Lat = min
	rt.NorthEast.Lat = max

	min = math.MaxFloat64
	max = -math.MaxFloat64

	for _, v := range this.Points {
		min = math.Min(min, v.Lng)
		max = math.Max(max, v.Lng)
	}
	rt.SouthWest.Lng = min
	rt.NorthEast.Lng = max
	return
}

/**
 * 计算两点之间的距离,两点坐标必须为经纬度
 * @param {point1} Point 点对象
 * @param {point2} Point 点对象
 * @returns {Number} 两点之间距离，单位为米
 */
func GetDistance(point1, point2 *Point) float64 {
	point1.Lng = getLoop(point1.Lng, -180, 180)
	point1.Lat = getRange(point1.Lat, -74, 74)
	point2.Lng = getLoop(point2.Lng, -180, 180)
	point2.Lat = getRange(point2.Lat, -74, 74)

	var x1, x2, y1, y2 float64
	x1 = DegreeToRad(point1.Lng)
	y1 = DegreeToRad(point1.Lat)
	x2 = DegreeToRad(point2.Lng)
	y2 = DegreeToRad(point2.Lat)

	return EARTHRADIUS * math.Acos((math.Sin(y1)*math.Sin(y2) + math.Cos(y1)*math.Cos(y2)*math.Cos(x2-x1)))

}

/**
 * 将度转化为弧度
 * @param {degree} Number 度
 * @returns {Number} 弧度
 */
func DegreeToRad(degree float64) float64 {
	return math.Pi * degree / 180
}

/**
 * 将弧度转化为度
 * @param {radian} Number 弧度
 * @returns {Number} 度
 */
func RadToDegree(rad float64) float64 {
	return (180 * rad) / math.Pi
}

/**
 * 将v值限定在a,b之间，纬度使用
 */
func getRange(v, a, b float64) float64 {
	v = math.Max(v, a)
	v = math.Min(v, b)
	return v
}

/**
 * 将v值限定在a,b之间，经度使用
 */
func getLoop(v, a, b float64) float64 {
	for v > b {
		v -= b - a
	}
	for v < a {
		v += b - a
	}
	return v
}

/**
 * 判断点是否在矩形内
 * @param {Point} point 点对象
 * @param {Bounds} bounds 矩形边界对象
 * @returns {Boolean} 点在矩形内返回true,否则返回false
 */
func IsPointInRect(point *Point, rt *Rect) bool {
	sw := rt.SouthWest
	ne := rt.NorthEast
	return (point.Lng >= sw.Lng && point.Lng <= ne.Lng && point.Lat >= sw.Lat && point.Lat <= ne.Lat)
}

/**
 * 判断点是否在圆形内
 * @param {Point} point 点对象
 * @param {Circle} circle 圆形对象
 * @returns {Boolean} 点在圆形内返回true,否则返回false
 */
func IsPointInCircle(point *Point, circle *Circle) bool {
	c := circle.Center
	r := circle.Radius
	dis := GetDistance(point, &c)
	if dis <= r {
		return true
	}
	return false
}

/**
 * 判断点是否在折线上
 * @param {Point} point 点对象
 * @param {Polyline} polyline 折线对象
 * @returns {Boolean} 点在折线上返回true,否则返回false
 */
func IsPointOnPolyline(point *Point, polyline *Polyline) bool {
	lineBounds := polyline.GetBounds()
	if IsPointInRect(point, lineBounds) == false {
		return false
	}
	//判断点是否在线段上，设点为Q，线段为P1P2 ，
	//判断点Q在该线段上的依据是：( Q - P1 ) × ( P2 - P1 ) = 0，且 Q 在以 P1，P2为对角顶点的矩形内
	for i := 0; i < len(polyline.Points)-1; i++ {
		var curPt = polyline.Points[i]
		var nextPt = polyline.Points[i+1]
		//首先判断point是否在curPt和nextPt之间，即：此判断该点是否在该线段的外包矩形内
		if point.Lng >= math.Min(curPt.Lng, nextPt.Lng) && point.Lng <= math.Max(curPt.Lng, nextPt.Lng) &&
			point.Lat >= math.Min(curPt.Lat, nextPt.Lat) && point.Lat <= math.Max(curPt.Lat, nextPt.Lat) {
			//判断点是否在直线上公式
			var precision = (curPt.Lng-point.Lng)*(nextPt.Lat-point.Lat) -
				(nextPt.Lng-point.Lng)*(curPt.Lat-point.Lat)
			if precision < F_ZERO && precision > -F_ZERO { //实质判断是否接近0
				return true
			}
		}
	}
	return false
}

/**
 * 判断点是否多边形内
 * @param {Point} point 点对象
 * @param {Polyline} polygon 多边形对象
 * @returns {Boolean} 点在多边形内返回true,否则返回false
 */
func IsPointInPolygon(point *Point, polygon *Polyline) bool {
	//首先判断点是否在多边形的外包矩形内，如果在，则进一步判断，否则返回false
	lineBounds := polygon.GetBounds()
	if IsPointInRect(point, lineBounds) == false {
		return false
	}

	var pts = polygon.Points
	//下述代码来源：http://paulbourke.net/geometry/insidepoly/，进行了部分修改
	//基本思想是利用射线法，计算射线与多边形各边的交点，如果是偶数，则点在多边形外，否则
	//在多边形内。还会考虑一些特殊情况，如点在多边形顶点上，点在多边形边上等特殊情况。
	var N = len(pts)
	var boundOrVertex = true //如果点位于多边形的顶点或边上，也算做点在多边形内，直接返回true
	var intersectCount = 0   //cross points count of x
	var precision = F_ZERO   //浮点类型计算时候与0比较时候的容差
	var p1, p2 *Point        //neighbour bound vertices
	var p = point            //测试点

	p1 = pts[0] //left vertex

	for i := 1; i <= N; i++ { //check all rays
		if p.Equals(p1) {
			return boundOrVertex //p is an vertex
		}

		p2 = pts[i%N]                                                             //right vertex
		if p.Lat < math.Min(p1.Lat, p2.Lat) || p.Lat > math.Max(p1.Lat, p2.Lat) { //ray is outside of our interests
			p1 = p2
			continue //next ray left point
		}

		if p.Lat > math.Min(p1.Lat, p2.Lat) && p.Lat < math.Max(p1.Lat, p2.Lat) { //ray is crossing over by the algorithm (common part of)
			if p.Lng <= math.Max(p1.Lng, p2.Lng) { //x is before of ray
				if p1.Lat == p2.Lat && p.Lng >= math.Min(p1.Lng, p2.Lng) { //overlies on a horizontal ray
					return boundOrVertex
				}

				if p1.Lng == p2.Lng { //ray is vertical
					if p1.Lng == p.Lng { //overlies on a vertical ray
						return boundOrVertex
					} else { //before ray
						intersectCount++
					}
				} else { //cross point on the left side
					var xinters = (p.Lat-p1.Lat)*(p2.Lng-p1.Lng)/(p2.Lat-p1.Lat) + p1.Lng //cross point of Lng
					if math.Abs(p.Lng-xinters) < precision {                              //overlies on a ray
						return boundOrVertex
					}

					if p.Lng < xinters { //before ray
						intersectCount++
					}
				}
			}
		} else { //special case when ray is crossing through the vertex
			if p.Lat == p2.Lat && p.Lng <= p2.Lng { //p crossing over p2
				var p3 = pts[(i+1)%N]                                                       //next vertex
				if p.Lat >= math.Min(p1.Lat, p3.Lat) && p.Lat <= math.Max(p1.Lat, p3.Lat) { //p.Lat lies between p1.Lat & p3.Lat
					intersectCount++
				} else {
					intersectCount += 2
				}
			}
		}
		p1 = p2 //next ray left point
	}

	if intersectCount%2 == 0 { //偶数在多边形外
		return false
	} else { //奇数在多边形内
		return true
	}

	return false
}

/**
 * 计算折线或者点数组的长度
 * @param {Polyline|Array<Point>} polyline 折线对象或者点数组
 * @returns {Number} 折线或点数组对应的长度
 */
func GetPolylineDistance(polyline *Polyline) float64 {
	//将polyline统一为数组
	var pts = polyline.Points

	if len(pts) < 2 { //小于2个点，返回0
		return 0
	}

	//遍历所有线段将其相加，计算整条线段的长度
	var totalDis float64
	for i := 0; i < len(pts)-1; i++ {
		var curPt = pts[i]
		var nextPt = pts[i+1]
		var dis = GetDistance(curPt, nextPt)
		totalDis += dis
	}

	return totalDis
}

/**
 * 计算多边形面或点数组构建图形的面积,注意：坐标类型只能是经纬度，且不适合计算自相交多边形的面积
 * @param {Polygon|Array<Point>} polygon 多边形面对象或者点数组
 * @returns {Number} 多边形面或点数组构成图形的面积
 */

func GetPolygonArea(polygon *Polyline) float64 {
	var pts = polygon.Points

	if len(pts) < 3 { //小于3个顶点，不能构建面
		return 0
	}

	var totalArea float64 = 0 //初始化总面积
	var LowX float64 = 0.0
	var LowY float64 = 0.0
	var MiddleX float64 = 0.0
	var MiddleY float64 = 0.0
	var HighX float64 = 0.0
	var HighY float64 = 0.0
	var AM float64 = 0.0
	var BM float64 = 0.0
	var CM float64 = 0.0
	var AL float64 = 0.0
	var BL float64 = 0.0
	var CL float64 = 0.0
	var AH float64 = 0.0
	var BH float64 = 0.0
	var CH float64 = 0.0
	var CoefficientL float64 = 0.0
	var CoefficientH float64 = 0.0
	var ALtangent float64 = 0.0
	var BLtangent float64 = 0.0
	var CLtangent float64 = 0.0
	var AHtangent float64 = 0.0
	var BHtangent float64 = 0.0
	var CHtangent float64 = 0.0
	var ANormalLine float64 = 0.0
	var BNormalLine float64 = 0.0
	var CNormalLine float64 = 0.0
	var OrientationValue float64 = 0.0
	var AngleCos float64 = 0.0
	var Sum1 float64 = 0.0
	var Sum2 float64 = 0.0
	var Count2 float64 = 0
	var Count1 float64 = 0
	var Sum float64 = 0.0
	var Radius = EARTHRADIUS //6378137.0,WGS84椭球半径
	var Count = len(pts)
	for i := 0; i < Count; i++ {
		if i == 0 {
			LowX = pts[Count-1].Lng * math.Pi / 180
			LowY = pts[Count-1].Lat * math.Pi / 180
			MiddleX = pts[0].Lng * math.Pi / 180
			MiddleY = pts[0].Lat * math.Pi / 180
			HighX = pts[1].Lng * math.Pi / 180
			HighY = pts[1].Lat * math.Pi / 180
		} else if i == Count-1 {
			LowX = pts[Count-2].Lng * math.Pi / 180
			LowY = pts[Count-2].Lat * math.Pi / 180
			MiddleX = pts[Count-1].Lng * math.Pi / 180
			MiddleY = pts[Count-1].Lat * math.Pi / 180
			HighX = pts[0].Lng * math.Pi / 180
			HighY = pts[0].Lat * math.Pi / 180
		} else {
			LowX = pts[i-1].Lng * math.Pi / 180
			LowY = pts[i-1].Lat * math.Pi / 180
			MiddleX = pts[i].Lng * math.Pi / 180
			MiddleY = pts[i].Lat * math.Pi / 180
			HighX = pts[i+1].Lng * math.Pi / 180
			HighY = pts[i+1].Lat * math.Pi / 180
		}
		AM = math.Cos(MiddleY) * math.Cos(MiddleX)
		BM = math.Cos(MiddleY) * math.Sin(MiddleX)
		CM = math.Sin(MiddleY)
		AL = math.Cos(LowY) * math.Cos(LowX)
		BL = math.Cos(LowY) * math.Sin(LowX)
		CL = math.Sin(LowY)
		AH = math.Cos(HighY) * math.Cos(HighX)
		BH = math.Cos(HighY) * math.Sin(HighX)
		CH = math.Sin(HighY)
		CoefficientL = (AM*AM + BM*BM + CM*CM) / (AM*AL + BM*BL + CM*CL)
		CoefficientH = (AM*AM + BM*BM + CM*CM) / (AM*AH + BM*BH + CM*CH)
		ALtangent = CoefficientL*AL - AM
		BLtangent = CoefficientL*BL - BM
		CLtangent = CoefficientL*CL - CM
		AHtangent = CoefficientH*AH - AM
		BHtangent = CoefficientH*BH - BM
		CHtangent = CoefficientH*CH - CM
		AngleCos = (AHtangent*ALtangent + BHtangent*BLtangent + CHtangent*CLtangent) / (math.Sqrt(AHtangent*AHtangent+BHtangent*BHtangent+CHtangent*CHtangent) * math.Sqrt(ALtangent*ALtangent+BLtangent*BLtangent+CLtangent*CLtangent))
		AngleCos = math.Acos(AngleCos)
		ANormalLine = BHtangent*CLtangent - CHtangent*BLtangent
		BNormalLine = 0 - (AHtangent*CLtangent - CHtangent*ALtangent)
		CNormalLine = AHtangent*BLtangent - BHtangent*ALtangent
		if AM != 0 {
			OrientationValue = ANormalLine / AM
		} else if BM != 0 {
			OrientationValue = BNormalLine / BM
		} else {
			OrientationValue = CNormalLine / CM
		}
		if OrientationValue > 0 {
			Sum1 += AngleCos
			Count1++
		} else {
			Sum2 += AngleCos
			Count2++
		}
	}
	var tempSum1, tempSum2 float64
	tempSum1 = Sum1 + (2*math.Pi*Count2 - Sum2)
	tempSum2 = (2*math.Pi*Count1 - Sum1) + Sum2
	if Sum1 > Sum2 {
		if (tempSum1 - float64(Count-2)*math.Pi) < 1 {
			Sum = tempSum1
		} else {
			Sum = tempSum2
		}
	} else {
		if (tempSum2 - float64(Count-2)*math.Pi) < 1 {
			Sum = tempSum2
		} else {
			Sum = tempSum1
		}
	}
	totalArea = (Sum - float64(Count-2)*math.Pi) * Radius * Radius

	return totalArea //返回总面积
}

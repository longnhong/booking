package ultility

import (
	"math"
)

const kmtomiles = float64(0.621371192)
const earthRadius = float64(6371)

func Haversine(latFrom float64, lngFrom float64, latTo float64, lngTo float64) (distance float64) {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lngTo - lngFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance = earthRadius * c

	return
}

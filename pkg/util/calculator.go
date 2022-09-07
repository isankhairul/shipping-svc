package util

import "math"

const (
	earthRaidusKm = 6371 // radius of the earth in kilometers.
)

func ConvertDegToRad(d float64) float64 {
	return d * math.Pi / 180
}

func CalculateDistanceInKm(lat1, long1, lat2, long2 float64) float64 {
	lat1 = ConvertDegToRad(lat1)
	long1 = ConvertDegToRad(long1)
	lat2 = ConvertDegToRad(lat2)
	long2 = ConvertDegToRad(long2)

	dLat := lat2 - lat1
	dLong := long2 - long1

	//a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLong/2), 2)

	//c = 2 ⋅ atan2( √a, √(1−a) )
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	//d = R * c
	d := earthRaidusKm * c
	return d
}

func CalculateVolume(length, weight, height float64) float64 {
	return length * weight * height
}

// length in cm
// wight in cm
// height in cm
func CalculateVolumeWeightKg(length, weight, height float64) float64 {
	return CalculateVolume(length, weight, height) / 6000
}

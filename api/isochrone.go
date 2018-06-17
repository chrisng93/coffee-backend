package api

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/chrisng93/coffee-backend/models"
	"googlemaps.github.io/maps"
)

// AvgWalkingSpeedMilesPerHour is the average walking speed in miles per hour.
const AvgWalkingSpeedMilesPerHour = 5

// NumOfAngles is the number of angles at which we calculate an isochrone.
const NumOfAngles = 12

// Tolerance is the percentage error we allow when finding travel times for an isochrone.
const Tolerance = 0.05

// EarthRadiusMiles is Earth's radius in miles. Used for Haversine formula.
const EarthRadiusMiles = 3961

func zip(x, y []float64) [][]float64 {
	output := make([][]float64, len(x))
	for i := 0; i < len(x); i++ {
		output[i] = []float64{x[i], y[i]}
	}
	return output
}

func sumOfRadiusDifferences(radius0, radius1 []float64) float64 {
	// Create slice where each index has a slice with a radius0 value and a radius1 value.
	zippedRadiuses := zip(radius0, radius1)
	sum := float64(0)
	for _, radii := range zippedRadiuses {
		sum += radii[0] - radii[1]
	}
	return sum
}

func degreesToRadians(angle float64) float64 {
	return angle * math.Pi / 180
}

func radiansToDegrees(radians float64) float64 {
	return 180 * radians / math.Pi
}

// Calculate lat/lng of point of radius distance and given angle away from origin using the
// Haversine formula. Needed because Earth is round. Sorry Kyrie.
func calculateLatLng(origin *models.Coordinates, radius float64, angle float64) []float64 {
	bearing := degreesToRadians(angle)
	lat1 := degreesToRadians(origin.Latitude)
	lng1 := degreesToRadians(origin.Longitude)
	// Haversine formula.
	lat2Radians := math.Asin(math.Sin(lat1)*math.Cos(radius/EarthRadiusMiles) + math.Cos(lat1)*math.Sin(radius/EarthRadiusMiles)*math.Cos(bearing))
	lng2Radians := lng1 + math.Atan2(math.Sin(bearing)*math.Sin(radius/EarthRadiusMiles)*math.Cos(lat1), math.Cos(radius/EarthRadiusMiles)-math.Sin(lat1)*math.Sin(lat2Radians))
	lat2 := radiansToDegrees(lat2Radians)
	lng2 := radiansToDegrees(lng2Radians)
	return []float64{lat2, lng2}
}

func coordinatesToString(lat float64, lng float64) string {
	coordinatesSlice := []string{
		strconv.FormatFloat(lat, 'f', 6, 64),
		",",
		strconv.FormatFloat(lng, 'f', 6, 64),
	}
	return strings.Join(coordinatesSlice, "")
}

func getDistanceMatrixResponse(
	googleMapsClient *maps.Client,
	origin *models.Coordinates,
	iso [][]float64,
) ([]string, []float64, error) {
	var destinationSlice []string
	for _, latLng := range iso {
		destinationSlice = append(destinationSlice, coordinatesToString(latLng[0], latLng[1]))
	}

	req := &maps.DistanceMatrixRequest{
		Origins:      []string{coordinatesToString(origin.Latitude, origin.Longitude)},
		Destinations: []string{strings.Join(destinationSlice, "|")},
		Mode:         "ModeWalking",
	}
	resp, err := googleMapsClient.DistanceMatrix(context.Background(), req)
	if err != nil {
		return nil, nil, err
	}

	var durations []float64
	for _, row := range resp.Rows {
		for _, element := range row.Elements {
			if element.Status == "OK" {
				durations = append(durations, element.Duration.Seconds()/60)
			}
		}
	}
	return resp.DestinationAddresses, durations, nil
}

func calculateIsochrones(
	googleMapsClient *maps.Client,
	origin *models.Coordinates,
	walkingTimeMin int64,
) ([][]float64, error) {
	var radius0, radius1, radius2, angles, radiusMin, radiusMax []float64
	var addresses0 []string
	var iso [][]float64
	for i := 0; i < NumOfAngles; i++ {
		// The radius slices are used to ___
		radius0 = append(radius0, 0)
		radius1 = append(radius1, float64(walkingTimeMin)*float64(AvgWalkingSpeedMilesPerHour)/60)
		radius2 = append(radius2, 0)
		// angles is used to ___
		angles = append(angles, float64(i*(360/NumOfAngles)))
		// addresses0 used to ___
		addresses0 = append(addresses0, "")
		// radiusMin and radiusMax used to ___
		radiusMin = append(radiusMin, 0)
		radiusMax = append(radiusMax, float64(AvgWalkingSpeedMilesPerHour)/60*float64(walkingTimeMin))
		// iso used to ___
		iso = append(iso, []float64{0, 0})
	}

	j := 0
	for sumOfRadiusDifferences(radius0, radius1) != 0 {
		for i := 0; i < NumOfAngles; i++ {
			iso[i] = calculateLatLng(origin, radius1[i], angles[i])
		}
		addresses, durations, err := getDistanceMatrixResponse(googleMapsClient, origin, iso)
		if err != nil {
			return nil, err
		}
		if len(addresses) == 0 || len(durations) == 0 {
			return [][]float64{}, nil
		}
		for i := 0; i < NumOfAngles; i++ {
			if durations[i] < (float64(walkingTimeMin)-float64(walkingTimeMin)*Tolerance) &&
				addresses0[i] != addresses[i] {
				radius2[i] = (radiusMax[i] + radius1[i]) / 2
				radiusMin[i] = radius1[i]
			} else if durations[i] > (float64(walkingTimeMin)+float64(walkingTimeMin)*Tolerance) &&
				addresses0[i] != addresses[i] {
				radius2[i] = (radiusMin[i] + radius1[i]) / 2
				radiusMax[i] = radius1[i]
			} else {
				radius2[i] = radius1[i]
			}
			addresses0[i] = addresses[i]
		}
		radius0 = radius1
		radius1 = radius2
		j++
		if j > 30 {
			return nil, errors.New("Isochrone calculation taking too long")
		}
	}
	return iso, nil
}

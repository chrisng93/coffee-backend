package api

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chrisng93/coffee-backend/models"
	"googlemaps.github.io/maps"
)

// MaxTries is the maximum number of isochrone calculating cycles before we error out.
const MaxTries = 30

// AvgWalkingSpeedMilesPerHour is the average walking speed in miles per hour.
const AvgWalkingSpeedMilesPerHour = 5

// FastestWalkingSpeedMilesPerHour is the upper limit on walking speed for humans in miles per
// hour. The Guinness World Record is 9.17 miles per hour.
const FastestWalkingSpeedMilesPerHour = 8

// NumOfAngles is the number of angles at which we calculate an isochrone.
const NumOfAngles = 12

// Tolerance is the percentage error we allow when finding travel times for an isochrone.
const Tolerance = 0.10

// EarthRadiusMiles is Earth's radius in miles. Used for Haversine formula.
const EarthRadiusMiles = 3961

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

// Turn lat/lng coordinates into a string representation for the request to the Google Maps
// Distance Matrix API.
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
) ([]float64, error) {
	var destinationSlice []string
	for _, latLng := range iso {
		destinationSlice = append(destinationSlice, coordinatesToString(latLng[0], latLng[1]))
	}

	req := &maps.DistanceMatrixRequest{
		Origins:      []string{coordinatesToString(origin.Latitude, origin.Longitude)},
		Destinations: []string{strings.Join(destinationSlice, "|")},
		Mode:         maps.TravelModeWalking,
	}
	resp, err := googleMapsClient.DistanceMatrix(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var durations []float64
	for _, row := range resp.Rows {
		for _, element := range row.Elements {
			if element.Status == "OK" {
				durations = append(durations, element.Duration.Seconds()/60)
			}
		}
	}
	return durations, nil
}

func calculateIsochrones(
	googleMapsClient *maps.Client,
	origin *models.Coordinates,
	walkingTimeMin int64,
) ([][]float64, error) {
	var radius0, radius1, angles, radiusMin, radiusMax, durations []float64
	var isochrone [][]float64
	for i := 0; i < NumOfAngles; i++ {
		// Radius slices are used to hold the radius distance from the origin for each angle that
		// we calculate isochrones for. These values are updated on each incremental calculation
		// of isochrones if we do not have a walking distance that's within the tolerance.
		radius0 = append(radius0, 0)
		radius1 = append(radius1, float64(walkingTimeMin)*float64(AvgWalkingSpeedMilesPerHour)/60)
		// angles is used to hold the respective angle for each angle in NumOfAngles. These
		// values are unchanged in the isochrone calculation process.
		angles = append(angles, float64(i*(360/NumOfAngles)))
		// radiusMin and radiusMax are initially the minimum and maximum distance for the route.
		// They're narrowed down in the isochrone calculation process.
		radiusMin = append(radiusMin, 0)
		radiusMax = append(radiusMax, float64(FastestWalkingSpeedMilesPerHour)/60*float64(walkingTimeMin))
		// durations is used to hold the walking durations from the Google Maps API for the
		// respective radius.
		durations = append(durations, 0)
		// isochrone holds the lat/lng of the radii.
		isochrone = append(isochrone, []float64{0, 0})
	}

	j := 0
	radiusDiffs := NumOfAngles
	// Allow for one radius diff to be greater than tolerance to lower response time. Throw this
	// value (if any) out later.
	for radiusDiffs > 1 {
		var tempRadius []float64
		radiusDiffs = 0
		// Use Haversine formula to calculate lat/lng for radius/angles from origin.
		for i := 0; i < NumOfAngles; i++ {
			isochrone[i] = calculateLatLng(origin, radius1[i], angles[i])
		}
		var err error
		// Call Google Maps Distance Matrix API to get the actual walking distance from the origin
		// for each of the lat/lngs calculated above.
		durations, err = getDistanceMatrixResponse(googleMapsClient, origin, isochrone)
		fmt.Println("found durations", durations)
		if err != nil {
			// TODO: Different message if over Google Maps API quota.
			return nil, err
		}
		if len(durations) == 0 {
			return [][]float64{}, nil
		}
		for i := 0; i < NumOfAngles; i++ {
			fmt.Println("checking duration", durations[i])
			if durations[i] < (float64(walkingTimeMin) - float64(walkingTimeMin)*Tolerance) {
				// Actual walking time below the tolerance. Increase the radius.
				tempRadius = append(tempRadius, (radiusMax[i]+radius1[i])/2)
				radiusMin[i] = radius1[i]
				radiusDiffs++
			} else if durations[i] > (float64(walkingTimeMin) + float64(walkingTimeMin)*Tolerance) {
				// Actual walking time above the tolerance. Reduce the radius.
				tempRadius = append(tempRadius, (radiusMin[i]+radius1[i])/2)
				radiusMax[i] = radius1[i]
				radiusDiffs++
			} else {
				tempRadius = append(tempRadius, radius1[i])
			}
		}
		radius0 = radius1
		radius1 = tempRadius
		j++
		if j > MaxTries {
			return nil, errors.New("Isochrone calculation taking too long")
		}
	}
	if radiusDiffs > 0 {
		return filterIsochrones(walkingTimeMin, isochrone, durations), nil
	}
	return isochrone, nil
}

// Filter out isochrones for any lat/lngs that aren't within tolerance. This is needed because we
// allow up to one value to be outside of the tolerance range to speed up response times.
func filterIsochrones(walkingTimeMin int64, isochrone [][]float64, durations []float64) [][]float64 {
	var filteredIsochrone [][]float64
	for i, duration := range durations {
		if duration >= float64(walkingTimeMin)-float64(walkingTimeMin)*Tolerance &&
			duration <= float64(walkingTimeMin)+float64(walkingTimeMin)*Tolerance {
			filteredIsochrone = append(filteredIsochrone, isochrone[i])
		}
	}
	return filteredIsochrone
}

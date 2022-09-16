// Package gameconstants provides various interesting constants for games.
package gameconstants

const (
	// Approximate movement speeds for the average person.
	WalkingSpeed = 1.4 // m/s
	RunningSpeed = 5.0 // m/s

	// Approximate view distance for the average person in meters for a
	// normal sized object.
	// NOTE: The earth's curvature limits the view distance to 5km
	ViewDistanceNormal = 1000.0 // m
	ViewDistanceDark   = 100.0  // m

	// Average view cone for the average person in degrees.
	ViewConeCentral     = 60.0  // degrees
	ViewConePeripherial = 180.0 // degrees

	// Approximate sensory distance for the average person in meters.
	DistanceHearing = 100.0 // m
	DistanceSmell   = 10.0  // m
	DistanceTouch   = 0.7   // m
	DistanceTaste   = 0.1   // m

	// Human carrying capacity.
	CarryingCapacity = 100.0 // kg

	// Horse properties.
	HorseTrotSpeed   = 5.0  // m/s
	HorseCanterSpeed = 8.0  // m/s
	HorseGallopSpeed = 12.0 // m/s
	// This is usually up to 20% of the horse's weight.
	// See: https://www.deephollowranch.com/how-much-weight-can-a-horse-carry/
	HorseCarryingCapacity = 120 // kg

	// Cart properties.
	// NOTE: Horse carts are twice as fast as oxen carts.
	CartSpeedEmpty       = 1.8 // m/s
	CartSpeedLaden       = 1.4 // m/s
	CartCarryingCapacity = 200 // kg (per horse)

	// Bird properties
	BirdSpeed            = 10.0 // m/s
	BirdCarryingCapacity = 0.2  // kg

	// Earth properties.
	EarthGravity                     = 9.81        // m/s^2
	EarthRadius                      = 6371.0      // km
	EarthSurface                     = 510100000.0 // km^2
	EarthMaxElevation                = 8848.0      // m above sea level
	EarthElevationTemperatureFalloff = 0.0065      // °C/m
	// EarthMinTemperature           = -89.2       // °C (extreme)
	// EarthMaxTemperature           = 56.7        // °C (extreme)
)

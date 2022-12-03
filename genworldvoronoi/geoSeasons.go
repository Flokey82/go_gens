package genworldvoronoi

import (
	"math"
)

// The seasons of the year change the day/night cycle as a sine wave,
// which has an effect on the day/night temperature.
// The extremes are at the poles, where days become almost 24 hours long or
// nitghts become almost 24 hours long.
//
// At the equator, seasons are not really noticeable, as the day/night cycle
// is constant the entire year. (so the amplitude of the sine wave is 0)
//
// There is a dry and a wet season instead of the 4 seasons, which is I think
// due to the northern and southern hemispheres switching seasons and the
// global winds might push humidity or dryness across the equator. (???)
//
// The day length is almost a square wave at the poles. I think this is a
// sine wave with amplitude cap... Like a distortion guitar effect.
// The amplitude at the equator is 0, and at the poles... let's say 10, which
// would be square enough when capped at 1.
//
// We should start by calculating the daily average temperature for each region over
// a year. Then we can think about day/night temperature differences...
//
// Given this information we are able to adapt plants and animals to the seasons,
// cultures, and so on.
// http://www.atmo.arizona.edu/students/courselinks/fall16/atmo336/lectures/sec4/seasons.html
// https://github.com/woodcrafty/PyETo/blob/0b7ac9f149f4c89c5b5759a875010c521aa07f0f/pyeto/fao.py#L198 !!!
// https://github.com/willbeason/worldproc/blob/28fd3f0188082ade001110a6a73edda4b987ccdd/pkg/climate/temperature.go

func (m *Geo) calcSolarRadiation(dayOfYear int) []float64 {
	res := make([]float64, m.mesh.numRegions)
	for i := range res {
		if math.Abs(m.LatLon[i][0]) > 90 {
			panic(m.LatLon[i][0])
		}
		latRad := degToRad(m.LatLon[i][0])
		res[i] = calcSolarRadiation(latRad, dayOfYear)
	}
	return res
}

// Calculate incoming solar (or shortwave) radiation, *Rs* (radiation hitting
// a horizontal plane after scattering by the atmosphere) from latitude, and
// day of year.
//
// 'latitude': Latitude [radians].
// 'dayOfYear': Day of year integer between 1 and 365 or 366).
//
// Returns incoming solar (or shortwave) radiation [MJ m-2 day-1]
func calcSolarRadiation(latRad float64, dayOfYear int) float64 {
	daylightHours := calcDaylightHoursByLatitudeAndDayOfYear(latRad, dayOfYear)
	sunshineHours := daylightHours * 0.7 // 70% of daylight hours

	sd := solarDeclination(dayOfYear)
	sha := sunsetHourAngle(latRad, sd)
	ird := invRelDistEarthSun(dayOfYear)

	// TODO: Use clearSkyRadiation to calculate et at a given altitude.
	et := extraterrRadiation(latRad, sd, sha, ird)
	sr := solRadFromSunHours(daylightHours, sunshineHours, et)

	// At the poles, we spread the solar radiation over a wider area since
	// the angle of the sun is really flat. I just guessed that this will
	// be following 1-sin(lat-solar declination) curve.
	return sr * (1 - math.Sin(math.Abs(latRad-sd)))
}

// Calculate daylight hours from latitude and day of year.
// Based on FAO equation 34 in Allen et al (1998).
//
// 'latitude': Latitude [radians]
// 'dayOfYear': Day of year integer between 1 and 365 or 366).
//
// Returns daylight hours.
func calcDaylightHoursByLatitudeAndDayOfYear(latRad float64, dayOfYear int) float64 {
	sd := solarDeclination(dayOfYear)
	sha := sunsetHourAngle(latRad, sd)
	return daylightHours(sha)
}

// Calculate incoming solar (or shortwave) radiation, *Rs* (radiation hitting
// a horizontal plane after scattering by the atmosphere) from relative
// sunshine duration.
// If measured radiation data are not available this method is preferable
// to calculating solar radiation from temperature. If a monthly mean is
// required then divide the monthly number of sunshine hours by number of
// days in the month and ensure that *et_rad* and *daylight_hours* was
// calculated using the day of the year that corresponds to the middle of
// the month. Based on equations 34 and 35 in Allen et al (1998).
//
// 'daylightHours': Number of daylight hours [hours].
// 'sunshineHours': Sunshine duration [hours].
// 'etRad': Extraterrestrial radiation [MJ m-2 day-1].
//
// Returns incoming solar (or shortwave) radiation [MJ m-2 day-1]
func solRadFromSunHours(daylightHours, sunshineHours, etRad float64) float64 {
	// 0.5 and 0.25 are default values of regression constants (Angstrom values)
	// recommended by FAO when calibrated values are unavailable.
	epsilon := 1e-13
	return (0.5*(sunshineHours+epsilon)/(daylightHours+epsilon) + 0.25) * etRad
}

// Solar constant [ MJ m-2 min-1]
const solarConstant = 0.0820

// Calculate sunset hour angle (*Ws*) from latitude and solar
// declination. Based on FAO equation 25 in Allen et al (1998).
//
// 'latitude': Latitude [radians].
// Note: *latitude* should be negative if it in the southern
// hemisphere, positive if in the northern hemisphere.
// 'solDec': Solar declination [radians].
//
// Returns sunset hour angle [radians].
func sunsetHourAngle(latRad float64, solDec float64) float64 {
	cos_sha := -math.Tan(latRad) * math.Tan(solDec)
	// If tmp is >= 1 there is no sunset, i.e. 24 hours of daylight
	// If tmp is <= 1 there is no sunrise, i.e. 24 hours of darkness
	// See http://www.itacanet.org/the-sun-as-a-source-of-energy/
	// part-3-calculating-solar-angles/
	// Domain of acos is -1 <= x <= 1 radians (this is not mentioned in FAO-56!)
	return math.Acos(math.Min(math.Max(cos_sha, -1.0), 1.0))
}

// Calculate solar declination from day of the year.
// Based on FAO equation 24 in Allen et al (1998).
//
// 'dayOfYear': Day of year integer between 1 and 365 or 366).
//
// Returns solar declination [radians]
func solarDeclination(dayOfYear int) float64 {
	return 0.409 * math.Sin((2.0*math.Pi/365.0)*float64(dayOfYear)-1.39)
}

// Calculate daylight hours from sunset hour angle.
// Based on FAO equation 34 in Allen et al (1998).
//
// 'sha': Sunset hour angle [rad].
//
// Returns daylight hours.
func daylightHours(sha float64) float64 {
	return (24.0 / math.Pi) * sha
}

// Estimate daily extraterrestrial radiation (*Ra*, 'top of the atmosphere
// radiation').
// Based on equation 21 in Allen et al (1998). If monthly mean radiation is
// required make sure *sol_dec*. *sha* and *irl* have been calculated using
// the day of the year that corresponds to the middle of the month.
// **Note**: From Allen et al (1998): "For the winter months in latitudes
// greater than 55 degrees (N or S), the equations have limited validity.
// Reference should be made to the Smithsonian Tables to assess possible
// deviations."
//
// 'latitude': Latitude [radians]
// 'solDec': Solar declination [radians].
// 'sha': Sunset hour angle [radians].
// 'ird': Inverse relative distance earth-sun [dimensionless].
//
// Returns daily extraterrestrial radiation [MJ m-2 day-1]
func extraterrRadiation(latitude, solDec, sha, ird float64) float64 {
	tmp1 := (24.0 * 60.0) / math.Pi
	tmp2 := sha * math.Sin(latitude) * math.Sin(solDec)
	tmp3 := math.Cos(latitude) * math.Cos(solDec) * math.Sin(sha)
	return tmp1 * solarConstant * ird * (tmp2 + tmp3)
}

// Estimate clear sky radiation from altitude and extraterrestrial radiation.
// Based on equation 37 in Allen et al (1998) which is recommended when
// calibrated Angstrom values are not available.
//
// 'altitude': Elevation above sea level [m]
// 'etRad': Extraterrestrial radiation [MJ m-2 day-1].
//
// Returns clear sky radiation [MJ m-2 day-1]
func clearSkyRadiation(altitude float64, etRad float64) float64 {
	return (0.00002*altitude + 0.75) * etRad
}

// Calculate the inverse relative distance between earth and sun from
// day of the year. Based on FAO equation 23 in Allen et al (1998).
//
// 'dayOfYear': Day of the year [1 to 366]
//
// Returns inverse relative distance between earth and the sun.
func invRelDistEarthSun(dayOfYear int) float64 {
	return 1 + (0.033 * math.Cos((2.0*math.Pi/365.0)*float64(dayOfYear)))
}

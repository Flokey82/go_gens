package gameconstants

const (
	// Land measurements.
	HectareToM2 = 10000.0        // hectars to m^2
	AcreToM2    = 4046.86        // acre to m^2
	AreToM2     = 100.0          // are to m^2
	SqInchToM2  = 0.00064516     // square inch to m^2
	SqFootToM2  = 0.09290304     // square foot to m^2
	SqYardToM2  = 0.83612736     // square yard to m^2
	SqMileToM2  = 2589988.110336 // square mile to m^2

	// Distance measurements.
	InchToM       = 0.0254  // inch to m
	FootToM       = 0.3048  // feet to m
	YardToM       = 0.9144  // yard to m
	MileToM       = 1609.0  // miles to m
	NauticalMiToM = 1852.0  // nautical miles to m
	PoleToM       = 5.0292  // pole to m
	RodToM        = 5.0292  // rod to m
	FurlongToM    = 201.168 // furlong to m
	HandToM       = 0.1016  // hand to m
	LinkToM       = 0.2012  // link to m
	ChainToM      = 20.1168 // chain to m
	CableToM      = 185.2   // cable to m
	LeagueToM     = 4828.03 // league to m

	// Weight measurements.
	OunceToKg       = 0.0283495 // ounces to kg
	PoundToKg       = 0.453592  // pounds to kg
	StoneToKg       = 6.35029   // stone to kg
	USTonToKg       = 907.18474 // US ton to kg
	ImperialTonToKg = 1016.05   // imperial ton to kg
)

// KelvinToCelsius converts a temperature in Kelvin to Celsius.
func KelvinToCelsius(k float64) float64 {
	return k - 273.15
}

// FahrenheitToCelsius converts a temperature in Fahrenheit to Celsius.
func FahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

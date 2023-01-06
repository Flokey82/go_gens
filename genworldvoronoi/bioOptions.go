package genworldvoronoi

type BioLevelOption func(*BioLevel)

func BioLevelName(name string) BioLevelOption {
	return func(b *BioLevel) {
		b.Name = name
	}
}

// BioLevelTempRange sets the minimum and maximum temperature range for a species in °C.
func BioLevelTempRange(min, max float64) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.TempRange = [2]float64{min, max}
	}
}

// BioLevelHumidityRange sets the minimum and maximum humidity range for a species in dm.
func BioLevelHumidityRange(min, max float64) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.HumRange = [2]float64{min, max}
	}
}

// BioLevelElevRange sets the minimum and maximum elevation range for a species from -1.0 - 1.0.
func BioLevelElevRange(min, max float64) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.ElevRange = [2]float64{min, max}
	}
}

// BioLevelSteepRange sets the minimum and maximum steepness range for a species from 0.0 - 1.0.
func BioLevelSteepRange(min, max float64) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.SteepRange = [2]float64{min, max}
	}
}

// BioLevelEcosphere sets the ecosphere for a species.
func BioLevelEcosphere(ecosphere EcosphereType) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.Ecosphere = ecosphere
	}
}

// BioLevelLocomotion sets the locomotion type for a species.
func BioLevelLocomotion(locomotion Locomotion) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesProperties.Locomotion = locomotion
	}
}

// BioLevelDigestion sets the digestion type for a species.
func BioLevelDigestion(digestion DigestiveSystem) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesProperties.Digestion = digestion
	}
}

// BioLevelPreferredBiomes sets the preferred biomes for a species.
func BioLevelPreferredBiomes(biomes ...int) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.PreferredBiomes = biomes
	}
}

// BioLevelAppendHereditary appends hereditary features of a species.
func BioLevelAppendHereditary(heredity ...BioProperty) BioLevelOption {
	return func(b *BioLevel) {
		b.Hereditary = append(b.Hereditary, heredity...)
	}
}

// BioLevelSize sets the size of a species.
func BioLevelSize(size SpeciesSize) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesProperties.Size = size
	}
}

// BioLevelRainRange sets the minimum and maximum rain range for a species in dm.
func BioLevelRainRange(min, max float64) BioLevelOption {
	return func(b *BioLevel) {
		b.SpeciesTolerances.RainRange = [2]float64{min, max}
	}
}

package village

// compMaps returns the count of keys in 'a' missing or exceeding
// their counterpart in 'b' (but not the other way round).
//
// Example:
//
// a["itemA"] = 2
// a["itemB"] = 1
// a["itemC"] = 1
//
// b["itemA"] = 2
// b["itemC"] = 10
// 
// Result:
// res["itemB"] == 1
func compMaps(a, b map[string]int) map[string]int {
	res := make(map[string]int)
	for key, val := range a {
		if b[key] < val {
			res[key] = val - b[key]
		}
	}
	return res
}

// addToMap adds the map 'add' to 'a'.
func addToMap(a, add map[string]int) map[string]int {
	for key, val := range add {
			a[key] += val
	}
	return a
}

// addMaps combines two maps into a third map
func addMaps(a, b map[string]int) map[string]int {
	res := make(map[string]int)
	for key, val := range a {
		res[key] += val
	}
	for key, val := range b {
		res[key] += val
	}
	return res
}
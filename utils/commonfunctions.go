package utils


func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func indexOf(element string, data []string) (int) {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
 }

func RemoveElementFromSlice(series []string, value string) []string{
	i := indexOf(value, series)
	if i != -1 {
		series[i] = series[len(series)-1] // Copy last element to index i.
		series[len(series)-1] = ""   // Erase last element (write zero value).
		series = series[:len(series)-1]   // Truncate slice.
	}
	return series
}


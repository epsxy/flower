package utils

func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Array2DContains(a [][]string, b string) bool {
	for _, i := range a {
		if ArrayContains(i, b) {
			return true
		}
	}
	return false
}

func AppendWithoutDuplicates(a []string, b string) []string {
	if ArrayContains(a, b) {
		return a
	}
	a = append(a, b)
	return a
}

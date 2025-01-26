package arrayx

func Contains[T comparable](arr []T, child T) bool {
	for _, c := range arr {
		if child == c {
			return true
		}
	}
	return false
}

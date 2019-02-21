package utils

func Contains(slice []uint, elem uint) bool {
	for _, n := range slice {
		if elem == n {
			return true
		}
	}
	return false
}

package game

// AbsInt returns the absolute value of the given integer.
// If the input number x is negative, it returns -x, otherwise it returns x.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

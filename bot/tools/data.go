package tools

func IsStringInSlice(s string, slc []string) bool {
	for _, v := range slc {
		if v == s {
			return true
		}
	}

	return false
}

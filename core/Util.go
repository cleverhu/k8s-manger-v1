package core

func IsValidLabel(m1, m2 map[string]string) bool {
	for key := range m2 {
		if m2[key] != m1[key] {
			return false
		}
	}

	return true
}

package discovery

// copyMap creates a copy of a string map
func (s *Scanner) copyMap(original map[string]string) map[string]string {
	if original == nil {
		return make(map[string]string)
	}

	copy := make(map[string]string, len(original))
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

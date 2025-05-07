package domain

func subtractSlices(main []string, toRemove []string) []string {
	removeMap := make(map[string]struct{})
	for _, item := range toRemove {
		removeMap[item] = struct{}{}
	}
	result := []string{}
	for _, item := range main {
		if _, exists := removeMap[item]; !exists {
			result = append(result, item)
		}
	}
	return result
}
func uniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0, len(slice))
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
		}
	}
	for entry, _ := range keys {
		list = append(list, entry)
	}
	return list
}

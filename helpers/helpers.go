package helpers

func MergeStringMaps(maps ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, mp := range maps {
		for k, v := range mp {
			res[k] = v
		}
	}

	return res
}

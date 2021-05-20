package tool

func DeleteRepeatList(data []string) []string {
	m := make(map[string]int)
	for i := range data {
		if _, ok := m[data[i]]; ok {
			continue
		} else {
			m[data[i]] = 1
		}
	}
	var res []string
	for key := range m {
		res = append(res, key)
	}

	return res
}

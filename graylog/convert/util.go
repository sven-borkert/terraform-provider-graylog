package convert

func ListToMap(data []interface{}, key string) map[string]interface{} {
	m := make(map[string]interface{}, len(data))
	for _, d := range data {
		elem, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		a, ok := elem[key].(string)
		if !ok {
			continue
		}
		delete(elem, key)
		m[a] = elem
	}
	return m
}

func MapToList(data map[string]interface{}, key string) []interface{} {
	list := make([]interface{}, 0, len(data))
	for k, d := range data {
		elem, ok := d.(map[string]interface{})
		if !ok {
			continue
		}
		elem[key] = k
		list = append(list, elem)
	}
	return list
}

func InterfaceListToStringList(data []interface{}) []string {
	list := make([]string, 0, len(data))
	for _, a := range data {
		s, ok := a.(string)
		if !ok {
			continue
		}
		list = append(list, s)
	}
	return list
}

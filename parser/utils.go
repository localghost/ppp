package parser

// Copies entries from src map which are not present in dest map.
func merge(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := dest[k]; !ok {
			dest[k] = v
		}
	}
}

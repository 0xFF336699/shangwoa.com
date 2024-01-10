package map2

func GetString(v interface{}, d string) string {
	switch v.(type) {
	case string:
		return v.(string)
	}
	return d
}

func GetInt(v interface{}, d int) int {
	switch v.(type) {
	case int:
		return v.(int)
	}
	return d
}

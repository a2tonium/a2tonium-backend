package certificate_json_generator

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

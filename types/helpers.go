package types

func mapFromTypeSet(i interface{}) map[string]interface{} {
	return i.(map[string]interface{})
}

func mapFromTypeSetList(i []interface{}) map[string]interface{} {
	for _, v := range i {
		return mapFromTypeSet(v)
	}
	return make(map[string]interface{})
}
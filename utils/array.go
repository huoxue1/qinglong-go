package utils

func In[T comparable](data T, array []T) bool {
	for _, t := range array {
		if t == data {
			return true
		}
	}
	return false
}

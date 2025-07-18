package converters

func ListSafe[T any](data []T) []T {
	if len(data) == 0 {
		return make([]T, 0)
	}
	return data
}

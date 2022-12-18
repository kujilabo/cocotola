package factory

func Factory1[T any](f1 func() (T, error)) (T, error) {
	var result T

	v1, err := f1()
	if err != nil {
		return result, err
	}
	return v1, nil
}
func PrintSlice[T any](s []T) {
	for _, v := range s {
		print(v)
	}
}

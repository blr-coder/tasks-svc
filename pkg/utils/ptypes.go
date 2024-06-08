package utils

func PBool(v bool) *bool {
	return &v
}

func Pointer[T any](v T) *T {
	return &v
}

func Value[T any](v *T) T {
	var d T
	if v != nil {
		return *v
	}
	return d
}

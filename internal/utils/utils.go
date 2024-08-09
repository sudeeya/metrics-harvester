package utils

func Int64Ptr(i int64) *int64 {
	ptr := new(int64)
	*ptr = i
	return ptr
}

func Float64Ptr(f float64) *float64 {
	ptr := new(float64)
	*ptr = f
	return ptr
}

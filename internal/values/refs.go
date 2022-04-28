package values

func StringP(value string) *string {
	return &value
}

func IntP(value int) *int {
	return &value
}

func FloatP(value float64) *float64 {
	return &value
}

package validation

func ToFloat64(val any) float64 {
	switch i := val.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	default:
		return 0
	}
}

func ToInt(val any) int {
	switch i := val.(type) {
	case int:
		return i
	case int8:
		return int(i)
	case int16:
		return int(i)
	case int32:
		return int(i)
	case int64:
		return int(i)
	case uint:
		return int(i)
	case uint8:
		return int(i)
	case uint16:
		return int(i)
	case uint32:
		return int(i)
	case uint64:
		return int(i)
	case float32:
		return int(i)
	case float64:
		return int(i)
	default:
		return 0
	}
}

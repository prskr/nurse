package validation

import (
	"encoding/json"
	"fmt"
)

func ParseJSONInt(got any) (int, error) {
	switch in := got.(type) {
	case float32, float64:
		return ToInt(in), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return ToInt(got), nil
	case []byte:
		var val json.Number
		if err := json.Unmarshal(in, &val); err != nil {
			return 0, err
		}
		if i, err := val.Int64(); err != nil {
			return 0, err
		} else {
			return int(i), nil
		}
	case string:
		var val json.Number
		if err := json.Unmarshal([]byte(in), &val); err != nil {
			return 0, err
		}
		if i, err := val.Int64(); err != nil {
			return 0, err
		} else {
			return int(i), nil
		}
	default:
		return 0, fmt.Errorf("cannot convert value %v to int", got)
	}
}

func ParseJSONFloat(got any) (float64, error) {
	switch in := got.(type) {
	case float32:
		return float64(in), nil
	case float64:
		return in, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i := got.(int)
		return float64(i), nil
	case []byte:
		var val json.Number
		if err := json.Unmarshal(in, &val); err != nil {
			return 0, err
		}
		return val.Float64()
	case string:
		var val json.Number
		if err := json.Unmarshal([]byte(in), &val); err != nil {
			return 0, err
		}
		return val.Float64()
	default:
		return 0, fmt.Errorf("cannot convert value %v to float", got)
	}
}

func ParseJSONString(got any) (string, error) {
	switch in := got.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", in), nil
	case []byte:
		return string(in), nil
	case string:
		return in, nil
	default:
		return "", fmt.Errorf("cannot convert value %v to float", got)
	}
}

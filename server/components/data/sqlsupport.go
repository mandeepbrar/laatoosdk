package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type StringArray []string

func (l StringArray) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *StringArray) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

type IntArray []int32

func (l IntArray) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *IntArray) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

type LongArray []int64

func (l LongArray) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *LongArray) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

type BoolArray []float64

func (l BoolArray) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *BoolArray) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

type FloatArray []float64

func (l FloatArray) Value() (driver.Value, error) {
	bytes, err := json.Marshal(l)
	return string(bytes), err
}

func (l *FloatArray) Scan(input interface{}) error {
	switch value := input.(type) {
	case string:
		return json.Unmarshal([]byte(value), l)
	case []byte:
		return json.Unmarshal(value, l)
	default:
		return errors.New("not supported")
	}
}

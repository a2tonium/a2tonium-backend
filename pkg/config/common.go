package config

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type value struct {
	value  any
	exists bool
}

type Value interface {
	// LookupString returns the config value and true if the key exists;
	// otherwise, it returns empty and false.
	LookupString() (string, bool)
	// String returns the config value or empty if not set.
	// Use LookupString to differentiate missing vs empty values.
	String() string
	// LookupInt returns the config value and true if the key exists;
	// otherwise, returns zero and false.
	LookupInt() (int, bool)
	// Int returns the config value or zero if not set.
	// Use LookupInt to distinguish missing from zero values.
	Int() int
	// LookupInt32 returns the int32 value and true if set; false if missing.
	// Use Int32 to get the value without existence check.
	LookupInt32() (int32, bool)
	// Int32 returns the int32 value or zero if not set.
	// Use LookupInt32 to distinguish missing from zero values.
	Int32() int32
	// LookupInt64 returns the int64 value and true if set; false if missing.
	// Use Int64 to get the value without existence check.
	LookupInt64() (int64, bool)
	// Int64 returns the int64 value or zero if not set.
	// Use LookupInt64 to distinguish missing from zero values.
	Int64() int64
	// LookupBoolean returns the bool value and true if set; false if missing.
	// Use Boolean to get the value without existence check.
	LookupBoolean() (bool, bool)
	// Boolean returns the bool value or false if not set.
	// Use LookupBoolean to distinguish missing from false values.
	Boolean() bool
	// LookupDuration returns the duration value and true if set; false if missing.
	// Use Duration to get the value without existence check.
	LookupDuration() (time.Duration, bool)
	// Duration returns the duration value or zero if not set.
	// Use LookupDuration to distinguish missing from zero values.
	Duration() time.Duration
	// LookupFloat32 returns the float32 value and true if set; false if missing.
	// Use Float32 to get the value without existence check.
	LookupFloat32() (float32, bool)
	// Float32 returns the float32 value or zero if not set.
	// Use LookupFloat32 to distinguish missing from zero values.
	Float32() float32
	// LookupFloat64 returns the float64 value and true if set; false if missing.
	// Use Float64 to get the value without existence check.
	LookupFloat64() (float64, bool)
	// Float64 returns the float64 value or zero if not set.
	// Use LookupFloat64 to distinguish missing from zero values.
	Float64() float64
	// LookupStrings returns the []string value and true if set; false if missing.
	// Use Strings to get the value without existence check.
	LookupStrings() ([]string, bool)
	// Strings returns the []string value or empty slice if not set.
	// Use LookupStrings to distinguish missing from empty values.
	Strings() []string
	// LookupInts returns the []int value and true if set; false if missing.
	// Use Ints to get the value without existence check.
	LookupInts() ([]int, bool)
	// Ints returns the []int value or empty slice if not set.
	// Use LookupInts to distinguish missing from empty values.
	Ints() []int // LookupInts32 retrieves the value of the config variable.
	// LookupInts32 returns the []int32 value and true if set; false if missing.
	// Use Ints32 to get the value without existence check.
	LookupInts32() ([]int32, bool)
	// Ints32 returns the []int32 value or empty slice if not set.
	// Use LookupInts32 to distinguish missing from empty values.
	Ints32() []int32
	// LookupInts64 returns the []int64 value and true if set; false if missing.
	// Use Ints64 to get the value without existence check.
	LookupInts64() ([]int64, bool)
	// Ints64 returns the []int64 value or empty slice if not set.
	// Use LookupInts64 to distinguish missing from empty values.
	Ints64() []int64
	// LookupFloats32 returns the []float32 value and true if set; false if missing.
	// Use Floats32 to get the value without existence check.
	LookupFloats32() ([]float32, bool)
	// Floats32 returns the []float32 value or empty slice if not set.
	// Use LookupFloats32 to distinguish missing from empty values.
	Floats32() []float32
	// LookupFloats64 returns the []float64 value and true if set; false if missing.
	// Use Floats64 to get the value without existence check.
	LookupFloats64() ([]float64, bool)
	// Floats64 returns the []float64 value or empty slice if not set.
	// Use LookupFloats64 to distinguish missing from empty values.
	Floats64() []float64
	// LookupBooleans returns the []bool value and true if set; false if missing.
	// Use Booleans to get the value without existence check.
	LookupBooleans() ([]bool, bool)
	// Booleans returns the []bool value or empty slice if not set.
	// Use LookupBooleans to distinguish missing from empty values.
	Booleans() []bool
	// LookupDurations returns the []time.Duration value and true if set; false if missing.
	// Use Durations to get the value without existence check.
	LookupDurations() ([]time.Duration, bool)
	// Durations returns the []time.Duration value or empty slice if not set.
	// Use LookupDurations to distinguish missing from empty values.
	Durations() []time.Duration
	// LookupMap returns the map[string]interface{} value and true if set; false if missing.
	// Use Map to get the value without existence check.
	LookupMap() (map[string]interface{}, bool)
	// Map returns the map[string]interface{} value or empty map if not set.
	// Use LookupMap to distinguish missing from empty values.
	Map() map[string]interface{}
}

func (v value) LookupString() (string, bool) {
	if !v.exists {
		return "", false
	}
	s, ok := v.value.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func (v value) String() string {
	val, _ := v.LookupString()
	return val
}

func (v value) LookupInt() (int, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(int)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Int() int {
	val, _ := v.LookupInt()
	return val
}

func (v value) LookupInt32() (int32, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(int32)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Int32() int32 {
	val, _ := v.LookupInt32()
	return val
}

func (v value) LookupInt64() (int64, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(int64)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Int64() int64 {
	val, _ := v.LookupInt64()
	return val
}

func (v value) LookupBoolean() (bool, bool) {
	if !v.exists {
		return false, false
	}
	b, ok := v.value.(bool)
	if !ok {
		return false, false
	}
	return b, true
}

func (v value) Boolean() bool {
	val, _ := v.LookupBoolean()
	return val
}

func (v value) LookupFloat32() (float32, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(float32)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Float32() float32 {
	val, _ := v.LookupFloat32()
	return val
}

func (v value) LookupFloat64() (float64, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(float64)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Float64() float64 {
	val, _ := v.LookupFloat64()
	return val
}

func (v value) LookupStrings() ([]string, bool) {
	if !v.exists {
		return nil, false
	}
	var (
		slice []interface{}
		val   interface{}
		res   []string
		str   string
		ok    bool
	)
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]string, 0, len(slice))

	for _, val = range slice {
		str, ok = val.(string)
		if !ok {
			return nil, false
		}
		res = append(res, str)
	}
	return res, true
}

func (v value) Strings() []string {
	val, _ := v.LookupStrings()
	return val
}

func (v value) LookupInts() ([]int, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []int
		i     int
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]int, 0, len(slice))

	for _, val = range slice {
		i, ok = val.(int)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Ints() []int {
	val, _ := v.LookupInts()
	return val
}

func (v value) LookupInts32() ([]int32, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []int32
		i     int32
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]int32, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(int32)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Ints32() []int32 {
	val, _ := v.LookupInts32()
	return val
}

func (v value) LookupInts64() ([]int64, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []int64
		i     int64
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]int64, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(int64)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Ints64() []int64 {
	val, _ := v.LookupInts64()
	return val
}

func (v value) LookupFloats32() ([]float32, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []float32
		i     float32
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]float32, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(float32)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Floats32() []float32 {
	val, _ := v.LookupFloats32()
	return val
}

func (v value) LookupFloats64() ([]float64, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []float64
		i     float64
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]float64, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(float64)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Floats64() []float64 {
	val, _ := v.LookupFloats64()
	return val
}

func (v value) LookupBooleans() ([]bool, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []bool
		b, ok bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]bool, 0, len(slice))
	for _, val = range slice {
		b, ok = val.(bool)
		if !ok {
			return nil, false
		}
		res = append(res, b)
	}
	return res, true
}

func (v value) Booleans() []bool {
	val, _ := v.LookupBooleans()
	return val
}

func (v value) LookupDuration() (time.Duration, bool) {
	var str string
	var ok bool
	if !v.exists {
		return 0, false
	}
	str, ok = v.value.(string)
	if !ok {
		return 0, false
	}

	var (
		err      error
		duration time.Duration
	)
	if str[len(str)-1:] == "y" {
		years, err := extractYears(str)
		if err != nil {
			return 0, false
		}
		duration = time.Duration(years) * (time.Hour * 24 * 365)

	} else {
		duration, err = time.ParseDuration(str)
		if err != nil {
			return 0, false
		}
	}

	return duration, true
}

func extractYears(input string) (int, error) {
	re := regexp.MustCompile(`(\d+)y`)
	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return 0, fmt.Errorf("no match found")
	}

	years, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}

	return years, nil
}

func (v value) Duration() time.Duration {
	val, _ := v.LookupDuration()
	return val
}

func (v value) LookupDurations() ([]time.Duration, bool) {
	var (
		slice []interface{}
		val   interface{}
		str   string
		res   []time.Duration
		d     time.Duration
		ok    bool
		err   error
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]time.Duration, 0, len(slice))
	for _, val = range slice {
		str, ok = val.(string)
		if !ok {
			return nil, false
		}
		d, err = time.ParseDuration(str)
		if err != nil {
			return nil, false
		}
		res = append(res, d)
	}
	return res, true
}

func (v value) Durations() []time.Duration {
	val, _ := v.LookupDurations()
	return val
}

func (v value) LookupMap() (map[string]interface{}, bool) {
	var m map[string]interface{}
	var ok bool
	if !v.exists {
		return nil, false
	}
	m, ok = v.value.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return m, true
}

func (v value) Map() map[string]interface{} {
	ret, _ := v.LookupMap()
	return ret
}

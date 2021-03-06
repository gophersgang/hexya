// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package nbutils

import (
	"fmt"
	"math"
	"strconv"
)

// CastToInteger casts the given val to int64 if it is
// a number type. Returns an error otherwise
func CastToInteger(val interface{}) (int64, error) {
	var res int64
	switch value := val.(type) {
	case int64:
		res = value
	case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32, float64:
		res, _ = strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
	default:
		return 0, fmt.Errorf("Value %v cannot be casted to int64", val)
	}
	return res, nil
}

// CastToFloat casts the given val to float64 if it is
// a number type. Panics otherwise
func CastToFloat(val interface{}) (float64, error) {
	var res float64
	switch value := val.(type) {
	case float64:
		res = value
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%d", value), 64)
	default:
		return 0, fmt.Errorf("Value %v cannot be casted to float64", val)
	}
	return res, nil
}

// Digits holds precision and scale information for a float (numeric) type:
//   - The precision: the total number of digits
//   - The scale: the number of digits to the right of the decimal point
//     (PostgresSQL definitions)
type Digits struct {
	Precision int8
	Scale     int8
}

// Round rounds the given val to the given precision.
// This function uses the future Go 1.10 implementation
func Round(val float64, precision Digits) float64 {
	const (
		mask  = 0x7FF
		shift = 64 - 11 - 1
		bias  = 1023

		signMask = 1 << 63
		fracMask = (1 << shift) - 1
		halfMask = 1 << (shift - 1)
		one      = bias << shift
	)

	bits := math.Float64bits(val)
	e := uint(bits>>shift) & mask
	switch {
	case e < bias:
		// Round abs(x)<1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= one // +-1
		}
	case e < bias+shift:
		// Round any abs(x)>=1 containing a fractional component [0,1).
		e -= bias
		bits += halfMask >> e
		bits &^= fracMask >> e
	}
	return math.Float64frombits(bits)
}

// Round32 rounds the given val to the given precision.
// This function is just a wrapper for Round() casted to float32
func Round32(val float32, precision Digits) float32 {
	return float32(Round(float64(val), precision))
}

// Compare 'value1' and 'value2' after rounding them according to the
// given precision. The returned values are per the following table:
//
//    value1 > value2 : true, false
//    value1 == value2: false, true
//    value1 < value2 : false, false
//
// A value is considered lower/greater than another value
// if their rounded value is different. This is not the same as having a
// non-zero difference!
//
// Example: 1.432 and 1.431 are equal at 2 digits precision,
// so this method would return 0
// However 0.006 and 0.002 are considered different (this method returns 1)
// because they respectively round to 0.01 and 0.0, even though
// 0.006-0.002 = 0.004 which would be considered zero at 2 digits precision.
//
// Warning: IsZero(value1-value2) is not equivalent to
// Compare(value1,value2) == _, true, as the former will round after
// computing the difference, while the latter will round before, giving
// different results for e.g. 0.006 and 0.002 at 2 digits precision.
func Compare(value1, value2 float64, precision Digits) (greater, equal bool) {
	if Round(value1, precision) == Round(value2, precision) {
		equal = true
		return
	}
	if Round(value1, precision) > Round(value2, precision) {
		greater = true
		return
	}
	return
}

// Compare32 'value1' and 'value2' after rounding them according to the
// given precision. This function is just a wrapper for Compare() with float32 values
func Compare32(value1, value2 float32, precision Digits) (greater, equal bool) {
	greater, equal = Compare(float64(value1), float64(value2), precision)
	return
}

// IsZero returns true if 'value' is small enough to be treated as
// zero at the given precision .
//
// Warning: IsZero(value1-value2) is not equivalent to
// Compare(value1,value2) == _, true, as the former will round after
// computing the difference, while the latter will round before, giving
// different results for e.g. 0.006 and 0.002 at 2 digits precision.
func IsZero(value float64, precision Digits) bool {
	epsilon := float64(10 ^ (-precision.Scale))
	if math.Abs(Round(value, precision)) < epsilon {
		return true
	}
	return false
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	// given
	var val interface{}
	val = 10
	asInt, err := ToInt(val)
	assert.Nil(t, err)
	assert.Equal(t, 10, asInt)

	val = "10"
	asInt, err = ToInt(val)
	assert.Nil(t, err)
	assert.Equal(t, 10, asInt)

	val = 10.0
	asInt, err = ToInt(val)
	assert.Nil(t, err)
	assert.Equal(t, 10, asInt)

	val = 10.1
	asInt, err = ToInt(val)
	assert.ErrorContains(t, err, "not an int")
	assert.Equal(t, 0, asInt)

	val = "hi"
	asInt, err = ToInt(val)
	assert.ErrorContains(t, err, "not an int")
	assert.Equal(t, 0, asInt)
}

func TestToBool(t *testing.T) {
	// given
	var val interface{}
	val = false
	asBool, err := ToBool(val)
	assert.Nil(t, err)
	assert.Equal(t, false, asBool)

	val = "false"
	asBool, err = ToBool(val)
	assert.Nil(t, err)
	assert.Equal(t, false, asBool)

	val = "true"
	asBool, err = ToBool(val)
	assert.Nil(t, err)
	assert.Equal(t, true, asBool)

	val = "NotABool"
	asBool, err = ToBool(val)
	assert.ErrorContains(t, err, "not a bool")
}

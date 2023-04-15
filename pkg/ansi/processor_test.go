package ansi

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/expect1.txt
var expect1String string

//go:embed testdata/expect2.txt
var expect2String string

func TestProcessFile1(t *testing.T) {
	result, err := ProcessFile("./testdata/input1.txt")
	assert.NoError(t, err)
	assert.Equal(t, expect1String, result)
}

func TestProcessFile3(t *testing.T) {
	result, err := ProcessFile("./testdata/input2.txt")
	assert.NoError(t, err)
	assert.Equal(t, expect2String, result)
}

package errorx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorX_WithStack(t *testing.T) {
	var err error
	err1 := WithStack(err, &ErrorX{})
	assert.Equal(t, err, err1)
}

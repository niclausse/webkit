package utils

import (
	"testing"
)

func TestToCamel(t *testing.T) {
	s := "created_at"
	_s := ToLowerCamel(s)
	if _s != "createdAt" {
		t.Error(_s)
	} else {
		t.Log(_s)
	}
}

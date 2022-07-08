package typex

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"github.com/penglin1995/webkit/utils"
	"strings"
)

type JsonX struct {
	Bytes []byte
	Valid bool
}

func (j JsonX) Value() (driver.Value, error) {
	s := j.trim()
	if len(s) == 0 {
		return nil, nil
	}

	return s, nil
}

func (j *JsonX) Scan(value interface{}) (err error) {
	if value == nil {
		j.Bytes, j.Valid = nil, false
		return
	}

	switch vt := value.(type) {
	case []byte:
		j.Bytes = append(j.Bytes[0:0], vt...)
	case string:
		j.Bytes = append(j.Bytes[0:0], utils.StringToBytes(vt)...)
	default:
		return fmt.Errorf("typex: jsonx support []byte,string only")
	}

	return
}

func (j JsonX) MarshalJsonX() ([]byte, error) {
	if j.IsNull() {
		return nullBytes, nil
	}

	return j.Bytes, nil
}

func (j *JsonX) UnmarshalJsonX(data []byte) (err error) {
	if len(data) == 0 || bytes.Equal(data, nullBytes) {
		j.Valid = false
		return
	}

	j.Bytes, j.Valid = append(j.Bytes[0:0], data...), true
	return
}

func (j JsonX) Equals(j1 JsonX) bool {
	return bytes.Equal(j.Bytes, j1.Bytes)
}

func (j JsonX) IsZero() bool {
	s := j.trim()
	return len(s) == 0 || s == "[]" || s == "{}"
}

func (j JsonX) IsNull() bool {
	if !j.Valid || len(j.Bytes) == 0 || bytes.Equal(j.Bytes, nullBytes) {
		return true
	}

	return false
}

func (j JsonX) trim() string {
	if j.IsNull() {
		return ""
	}

	return strings.ReplaceAll(strings.ReplaceAll(utils.BytesToString(j.Bytes), "\n", ""), "\t", "")
}

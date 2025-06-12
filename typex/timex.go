package typex

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/niclausse/webkit/v2/utils"
	"strings"
	"time"
)

const (
	DefaultFormat = "2006-01-02 15:04:05"
	DateFormat    = "2006-01-02"
)

var nullBytes = []byte("null")
var emptyString = []byte("\"\"")

type TimeX struct {
	sql.NullTime
}

func Now() TimeX {
	return TimeX{NullTime: sql.NullTime{Time: time.Now(), Valid: true}}
}

// Scan implements the Scanner interface.
func (t *TimeX) Scan(value interface{}) error {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	switch vt := value.(type) {
	case time.Time:
		t.Time = vt
		t.Valid = true
	case int64:
		t.Time = time.Unix(vt, 0)
		t.Valid = true
	}

	return nil
}

// Value implements the driver Valuer interface.
func (t TimeX) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	return t.Time, nil
}

func (t TimeX) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullBytes, nil
	}

	return utils.StringToBytes(t.Time.Format(`"2006-01-02 15:04:05"`)), nil
}

func (t *TimeX) UnmarshalJSON(data []byte) (err error) {
	t.Valid = false

	if bytes.Equal(data, nullBytes) || bytes.Equal(data, emptyString) {
		return
	}

	text := tidy(data)

	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006",
		"2006-01",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"2006/01",
		"",
	}

	for _, layout := range layouts {
		t.Time, err = time.ParseInLocation(layout, text, time.Local)
		if err != nil {
			continue
		}

		t.Valid = true
		break
	}

	if !t.Valid {
		return fmt.Errorf("typex: unsupport time format: %w", err)
	}

	return
}

func (t TimeX) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(DefaultFormat)
}

func (t *TimeX) Date() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(DateFormat)
}

func tidy(data []byte) string {
	s := utils.BytesToString(data)

	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		return s[1 : len(s)-1]
	}

	return s
}

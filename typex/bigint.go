package typex

import (
	"bytes"
	"fmt"
	"github.com/niclausse/webkit/v2/utils"
	"strconv"
)

// Bigint 大型整数类型， 吐出string类型， 防止前端number超出上限报错
type Bigint int64

func (i Bigint) Int64() int64 {
	return int64(i)
}

func (i Bigint) MarshalJSON() ([]byte, error) {
	return utils.StringToBytes(fmt.Sprintf("\"%v\"", i)), nil
}

func (i *Bigint) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, "\"")
	bigint, err := string2id(utils.BytesToString(b))
	if err != nil {
		return err
	}

	*i = bigint
	return nil
}

func string2id(idStr string) (Bigint, error) {
	if len(idStr) == 0 {
		return Bigint(0), nil
	}

	i, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("typex: unable to cast %s to int64: %v", idStr, err)
	}

	return Bigint(i), nil
}

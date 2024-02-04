package exportx

import (
	"strconv"
	"testing"
)

func TestCSV_Export_StringSlice1(t *testing.T) {
	var (
		header = []string{"id", "姓名", "年龄", "性别", "城市"}
		data   = make([][]string, 0, len(header))
	)

	res, err := NewCSV().Export(data, WithFilename("人员统计"), WithHeader(header))
	if err != nil {
		t.Errorf("TestCSV_Export1 Failed: %+v", err)
		return
	}

	t.Logf("TestCSV_Export1 Success: %s", string(res))
}

func TestCSV_Export_StringSlice2(t *testing.T) {
	var (
		header = []string{"id", "姓名", "年龄", "性别", "城市"}
		data   = make([][]string, 0)
	)

	for i := 0; i < 100; i++ {
		d := make([]string, len(header), len(header))
		d[0] = strconv.Itoa(i + 1)
		d[1] = "姓名" + strconv.Itoa(i+1)
		d[2] = "年龄" + strconv.Itoa(i+1)
		d[3] = "性别" + strconv.Itoa(i+1)
		d[4] = "城市" + strconv.Itoa(i+1)
		data = append(data, d)
	}

	res, err := NewCSV().Export(data, WithFilename("人员统计"), WithHeader(header))
	if err != nil {
		t.Errorf("TestCSV_Export_StringSlice2 Failed: %+v", err)
		return
	}

	t.Logf("TestWriteCscStringSlice1 Success: %s", string(res))
}

func TestCSV_Export_StructSlice1(t *testing.T) {
	type User struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
		City   string `json:"city"`
	}

	users := make([]*User, 0, 100)
	for i := 0; i < cap(users); i++ {
		users = append(users, &User{
			ID:     i + 1,
			Name:   "姓名" + strconv.Itoa(i+1),
			Age:    i + 1,
			Gender: "男",
			City:   "城市" + strconv.Itoa(i+1),
		})
	}

	res, err := NewCSV().Export(users, WithFilename("struct人员统计"))
	if err != nil {
		t.Errorf("TestCSV_WriteStructSlice Failed: %+v", err)
		return
	}
	t.Logf("TestCSV_WriteStructSlice Success: %s", string(res))
}

func TestCSV_ExportLocal_StringSlice(t *testing.T) {
	var (
		header = []string{"id", "姓名", "年龄", "性别", "城市"}
		data   = make([][]string, 0)
	)

	for i := 0; i < 100; i++ {
		d := make([]string, len(header), len(header))
		d[0] = strconv.Itoa(i + 1)
		d[1] = "姓名" + strconv.Itoa(i+1)
		d[2] = "年龄" + strconv.Itoa(i+1)
		d[3] = "性别" + strconv.Itoa(i+1)
		d[4] = "城市" + strconv.Itoa(i+1)
		data = append(data, d)
	}

	res, err := NewCSV().ExportLocal(data, WithFilename("人员统计"), WithHeader(header))
	if err != nil {
		t.Errorf("TestCSV_Export_StringSlice2 Failed: %+v", err)
		return
	}

	t.Logf("TestWriteCscStringSlice1 Success: %s", res)
}

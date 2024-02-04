package exportx

import (
	"os"
	"strconv"
	"testing"
)

func TestExcel_Export_Normal_Map(t *testing.T) {
	data := make([]map[string]string, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, map[string]string{
			"id":      strconv.Itoa(i + 1),
			"name":    "测试" + strconv.Itoa(i+1),
			"age":     strconv.Itoa(i + 1),
			"gender":  "男",
			"address": "测试" + strconv.Itoa(i+1),
		})
	}

	ex := NewExcel()

	content, err := ex.Export(
		data,
		WithHeader([]string{"序号", "姓名", "年龄", "性别", "地址"}),
		WithHeaderKeys([]string{"id", "name", "age", "gender", "address"}),
	)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	file, _ := os.Create("./export_files/test1.xlsx")
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestExcel_Export_Normal_Struct(t *testing.T) {
	type User struct {
		Id      string `json:"id" excel:"id"`
		Name    string `json:"name" excel:"name"`
		Age     string `json:"age" excel:"age"`
		Gender  string `json:"gender" excel:"gender"`
		Address string `json:"address" excel:"address"`
	}

	data := make([]*User, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, &User{
			Id:      strconv.Itoa(i + 1),
			Name:    "测试" + strconv.Itoa(i+1),
			Age:     strconv.Itoa(i + 1),
			Gender:  "男",
			Address: "测试" + strconv.Itoa(i+1),
		})
	}

	ex := NewExcel()

	content, err := ex.Export(
		data,
		WithHeader([]string{"序号", "姓名", "年龄", "性别", "地址"}),
		WithHeaderKeys([]string{"id", "name", "age", "gender", "address"}),
	)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	file, _ := os.Create("./export_files/test2.xlsx")
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		t.Errorf("%+v", err)
		return
	}
}

func TestExcel_ExportLocal_Normal_Map(t *testing.T) {
	data := make([]map[string]string, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, map[string]string{
			"id":      strconv.Itoa(i + 1),
			"name":    "测试" + strconv.Itoa(i+1),
			"age":     strconv.Itoa(i + 1),
			"gender":  "男",
			"address": "测试" + strconv.Itoa(i+1),
		})
	}

	ex := NewExcel()

	path, err := ex.ExportLocal(
		data,
		WithHeader([]string{"序号", "姓名", "年龄", "性别", "地址"}), WithFilename("学员"),
		WithHeaderKeys([]string{"id", "name", "age", "gender", "address"}),
	)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	t.Logf("%+v", path)
}

func TestExcel_ExportLocal_Normal_Struct(t *testing.T) {
	type User struct {
		Id      string `json:"id" excel:"id"`
		Name    string `json:"name" excel:"姓名"`
		Age     string `json:"age" excel:"age"`
		Gender  string `json:"gender" excel:"gender"`
		Address string `json:"address" excel:"address"`
	}

	data := make([]*User, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, &User{
			Id:      strconv.Itoa(i + 1),
			Name:    "测试" + strconv.Itoa(i+1),
			Age:     strconv.Itoa(i + 1),
			Gender:  "男",
			Address: "测试" + strconv.Itoa(i+1),
		})
	}

	ex := NewExcel(WithExportType(ExcelExportTypeNormal))

	content, err := ex.ExportLocal(data, WithFilename("学员"))
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	t.Logf("%s", content)
}

func TestExcel_ExportLocal_Stream_Map(t *testing.T) {
	data := make([]map[string]string, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, map[string]string{
			"id":      strconv.Itoa(i + 1),
			"name":    "测试" + strconv.Itoa(i+1),
			"age":     strconv.Itoa(i + 1),
			"gender":  "男",
			"address": "测试" + strconv.Itoa(i+1),
		})
	}

	ex := NewExcel(WithExportType(ExcelExportTypeStream))

	path, err := ex.ExportLocal(data, WithHeader([]string{"序号", "姓名", "年龄", "性别", "地址"}), WithHeaderKeys([]string{"id", "name", "age", "gender", "address"}), WithFilename("学员"))
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	t.Logf("%+v", path)
}

func TestExcel_ExportLocal_Stream_Struct(t *testing.T) {
	type User struct {
		Id      string `json:"id" excel:"id"`
		Name    string `json:"name" excel:"name"`
		Age     string `json:"age" excel:"age"`
		Gender  string `json:"gender" excel:"gender"`
		Address string `json:"address" excel:"address"`
	}

	data := make([]*User, 0, 100)
	for i := 0; i < 100; i++ {
		data = append(data, &User{
			Id:      strconv.Itoa(i + 1),
			Name:    "测试" + strconv.Itoa(i+1),
			Age:     strconv.Itoa(i + 1),
			Gender:  "男",
			Address: "测试" + strconv.Itoa(i+1),
		})
	}

	ex := NewExcel(WithExportType(ExcelExportTypeStream))

	content, err := ex.ExportLocal(data, WithFilename("学员"))
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	t.Logf("%s", content)
}

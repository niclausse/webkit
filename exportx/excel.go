package exportx

import (
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

type Excel struct {
	*config
	*ExcelConfig
	file  *excelize.File
	isPtr bool
}

type ExcelExportType string

const (
	GenExcelTypeNormal ExcelExportType = "normal"
	GenExcelTypeStream ExcelExportType = "stream"
)

type ExcelConfig struct {
	exportType ExcelExportType
	headerKeys []string // required for map data
	sheetName  string   // sheetName to replace default Sheet1
}

type ExcelOption func(opt *ExcelConfig)

func WithExportType(exportType ExcelExportType) ExcelOption {
	return func(c *ExcelConfig) {
		c.exportType = exportType
	}
}

func WithHeaderKeys(keys []string) ExcelOption {
	return func(c *ExcelConfig) {
		c.headerKeys = keys
	}
}

func WithSheetName(sheetName string) ExcelOption {
	return func(c *ExcelConfig) {
		c.sheetName = sheetName
	}
}

func NewExcel(opts ...ExcelOption) Exporter {
	ex := &Excel{
		config:      &config{},
		ExcelConfig: &ExcelConfig{},
	}
	for _, opt := range opts {
		opt(ex.ExcelConfig)
	}

	return ex
}

func (e *Excel) Export(data interface{}, opts ...Option) ([]byte, error) {
	for _, opt := range opts {
		opt(e.config)
	}

	if len(e.header) == 0 {
		return nil, ErrInvalidHeader
	}

	if len(e.headerKeys) != len(e.header) {
		return nil, ErrInvalidHeaderKeys
	}

	var (
		err error
	)

	if e.exportType == GenExcelTypeStream {
		err = e.stream(data)
	} else {
		err = e.normal(data)
	}
	if err != nil {
		return nil, err
	}

	bf, err := e.file.WriteToBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to write to buffer")
	}

	return bf.Bytes(), nil
}

func (e *Excel) ExportLocal(data interface{}, opts ...Option) (string, error) {
	for _, opt := range opts {
		opt(e.config)
	}

	if len(e.header) == 0 {
		return "", ErrInvalidHeader
	}

	if len(e.headerKeys) != len(e.header) {
		return "", ErrInvalidHeaderKeys
	}

	var (
		err error
	)

	if e.exportType == GenExcelTypeStream {
		err = e.stream(data)
	} else {
		err = e.normal(data)
	}
	if err != nil {
		return "", err
	}

	if err = e.initSaveProfile(ExcelExt); err != nil {
		return "", err
	}

	return e.savePath, e.file.SaveAs(e.savePath)
}

func (e *Excel) Close() error {
	return e.file.Close()
}

func (e *Excel) stream(data interface{}) (err error) {
	return // TODO
}

func (e *Excel) normal(data interface{}) (err error) {
	e.file = excelize.NewFile()

	if e.sheetName == "" {
		e.sheetName = "Sheet1"
	} else {
		if err = e.file.DeleteSheet("Sheet1"); err != nil {
			return errors.Wrap(err, "failed to delete default sheet")
		}
		if _, err = e.file.NewSheet(e.sheetName); err != nil {
			return errors.Wrap(err, "failed to create sheet")
		}
	}

	// try set header
	for i := 0; i < len(e.header); i++ {
		if err = e.file.SetCellValue(e.sheetName, transHeaderIndexToSheetPosition(i+1)+"1", e.header[i]); err != nil {
			return errors.Wrap(err, "failed to set header value")
		}
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return ErrInvalidDataFormat
	}

	typ := rv.Type().Elem()

	switch typ.Kind() {
	case reflect.Map:
		arr, ok := data.([]map[string]string)
		if !ok {
			return ErrInvalidDataFormat
		}

		if len(e.header) == 0 {
			return errors.New("exportx: header option is required for map data")
		}

		// set data
		for i := 0; i < len(arr); i++ {
			rowNum := i + 2 // 从第二行开始写数据(第一行写header)
			for j := 0; j < len(e.headerKeys); j++ {
				if err = e.file.SetCellValue(e.sheetName, transHeaderIndexToSheetPosition(j+1)+strconv.Itoa(rowNum), arr[i][e.headerKeys[j]]); err != nil {
					return errors.Wrap(err, "failed to set data value")
				}
			}
		}
	case reflect.Pointer:
		typ = typ.Elem()
		e.isPtr = true
		fallthrough
	case reflect.Struct:
		// set header
		if len(e.header) == 0 {
			for i := 0; i < typ.NumField(); i++ {
				tag, ok := typ.Field(i).Tag.Lookup("excel")
				if !ok {
					tag = typ.Field(i).Name
				}

				if err = e.file.SetCellValue(e.sheetName, transHeaderIndexToSheetPosition(i+1)+"1", tag); err != nil {
					return errors.Wrap(err, "failed to set header value")
				}
			}
		}

		// set data
		for i := 0; i < rv.Len(); i++ {
			rowNum := i + 2 // 从第二行开始写数据(第一行写header)
			fieldNum := typ.NumField()
			for j := 0; j < fieldNum; j++ {
				var ev reflect.Value
				if e.isPtr {
					ev = rv.Index(i).Elem().Field(j)
				} else {
					ev = rv.Index(i).Field(j)
				}

				if err = e.file.SetCellValue(e.sheetName, transHeaderIndexToSheetPosition(j+1)+strconv.Itoa(rowNum), ev); err != nil {
					return errors.Wrap(err, "failed to set data value")
				}
			}
		}
	default:
		return ErrInvalidDataFormat
	}

	return
}

// according to the header slice index to set sheet position
func transHeaderIndexToSheetPosition(index int) string {
	var (
		Str  string
		k    int
		temp []int // 保存转化后每一位数据的值，然后通过索引的方式匹配A-Z
	)
	// 用来匹配的字符A-Z
	Slice := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	if index > 26 { // 数据大于26需要进行拆分
		for {
			k = index % 26 // 从个位开始拆分，如果求余为0，说明末尾为26，也就是Z，如果是转化为26进制数，则末尾是可以为0的，这里必须为A-Z中的一个
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			index = (index - k) / 26 // 减去index最后一位数的值，因为已经记录在temp中
			if index <= 26 {         // 小于等于26直接进行匹配，不需要进行数据拆分
				temp = append(temp, index)
				break
			}
		}
	} else {
		return Slice[index]
	}
	for _, value := range temp {
		Str = Slice[value] + Str // 因为数据切分后存储顺序是反的，所以Str要放在后面
	}
	return Str
}

package exportx

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"log"
	"reflect"
)

type Excel struct {
	*config
	*ExcelConfig
	file         *excelize.File
	streamWriter *excelize.StreamWriter
	isPtr        bool
}

type ExcelExportType string

const (
	ExcelExportTypeNormal ExcelExportType = "normal"
	ExcelExportTypeStream ExcelExportType = "stream"
)

type ExcelConfig struct {
	exportType ExcelExportType
	sheetName  string // sheetName to replace default Sheet1
}

type ExcelOption func(opt *ExcelConfig)

func WithExportType(exportType ExcelExportType) ExcelOption {
	return func(c *ExcelConfig) {
		c.exportType = exportType
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
	defer e.close()

	for _, opt := range opts {
		opt(e.config)
	}

	var (
		err error
		bf  *bytes.Buffer
	)

	e.file = excelize.NewFile()

	if e.sheetName == "" {
		e.sheetName = "Sheet1"
	} else {
		if err = e.file.DeleteSheet("Sheet1"); err != nil {
			return nil, errors.Wrap(err, "exportx: failed to delete default sheet")
		}
		if _, err = e.file.NewSheet(e.sheetName); err != nil {
			return nil, errors.Wrap(err, "exportx: failed to create sheet")
		}
	}

	if e.exportType == ExcelExportTypeStream {
		e.streamWriter, err = e.file.NewStreamWriter(e.sheetName)
		if err != nil {
			return nil, errors.Wrap(err, "exportx: failed to create stream writer")
		}
		err = e.stream(data)
	} else {
		err = e.normal(data)
	}
	if err != nil {
		return nil, err
	}

	bf, err = e.file.WriteToBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to write to buffer")
	}

	return bf.Bytes(), nil
}

func (e *Excel) ExportLocal(data interface{}, opts ...Option) (string, error) {
	defer e.close()

	for _, opt := range opts {
		opt(e.config)
	}

	var err error

	e.file = excelize.NewFile()

	if e.sheetName == "" {
		e.sheetName = "Sheet1"
	} else {
		if err = e.file.DeleteSheet("Sheet1"); err != nil {
			return "", errors.Wrap(err, "exportx: failed to delete default sheet")
		}
		if _, err = e.file.NewSheet(e.sheetName); err != nil {
			return "", errors.Wrap(err, "exportx: failed to create sheet")
		}
	}

	if e.exportType == ExcelExportTypeStream {
		e.streamWriter, err = e.file.NewStreamWriter(e.sheetName)
		if err != nil {
			return "", errors.Wrap(err, "exportx: failed to create stream writer")
		}
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

func (e *Excel) close() {
	if e.streamWriter != nil {
		if err := e.streamWriter.Flush(); err != nil {
			log.Printf("exportx: failed to flush stream writer: %v", err)
		}
	}
	if e.file != nil {
		if err := e.file.Close(); err != nil {
			log.Printf("exportx: failed to close file: %v", err)
		}
	}
}

func (e *Excel) stream(data interface{}) (err error) {
	// try set header
	if len(e.header) > 0 {
		headers := make([]interface{}, 0, len(e.header))
		for i := 0; i < len(e.header); i++ {
			headers = append(headers, e.header[i])
		}
		if err = e.streamWriter.SetRow("A1", headers); err != nil {
			return errors.Wrap(err, "fail to set header row")
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

		if len(e.headerKeys) != len(e.header) {
			return errors.New("exportx: length of headerKeys is not equal to length of header")
		}

		// set data
		for i := 0; i < len(arr); i++ {
			rowNum := i + 2 // 从第二行开始写数据(第一行写header)
			rowValues := make([]interface{}, 0, len(e.headerKeys))
			for j := 0; j < len(e.headerKeys); j++ {
				rowValues = append(rowValues, arr[i][e.headerKeys[j]])
			}

			coordinate, _ := excelize.CoordinatesToCellName(1, rowNum)
			if err = e.streamWriter.SetRow(coordinate, rowValues); err != nil {
				return errors.Wrap(err, "failed to set data value")
			}
		}
	case reflect.Pointer:
		typ = typ.Elem()
		e.isPtr = true
		fallthrough
	case reflect.Struct:
		// set header
		if len(e.header) == 0 {
			headers := make([]interface{}, 0, typ.NumField())
			for i := 0; i < typ.NumField(); i++ {
				tag, ok := typ.Field(i).Tag.Lookup("excel")
				if !ok {
					tag = typ.Field(i).Name
				}
				headers = append(headers, tag)
			}

			if err = e.streamWriter.SetRow("A1", headers); err != nil {
				return errors.Wrap(err, "failed to set header value")
			}
		}

		// set data
		for i := 0; i < rv.Len(); i++ {
			rowNum := i + 2 // 从第二行开始写数据(第一行写header)
			rowValues := make([]interface{}, 0, typ.NumField())
			fieldNum := typ.NumField()
			for j := 0; j < fieldNum; j++ {
				var ev reflect.Value
				if e.isPtr {
					ev = rv.Index(i).Elem().Field(j)
				} else {
					ev = rv.Index(i).Field(j)
				}

				rowValues = append(rowValues, ev.Interface())
			}

			coordinate, _ := excelize.CoordinatesToCellName(1, rowNum)
			if err = e.streamWriter.SetRow(coordinate, rowValues); err != nil {
				return errors.Wrap(err, "failed to set data value")
			}
		}
	default:
		return ErrInvalidDataFormat
	}

	return
}

func (e *Excel) normal(data interface{}) (err error) {
	// try set header
	for i := 0; i < len(e.header); i++ {
		coordinate, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err = e.file.SetCellValue(e.sheetName, coordinate, e.header[i]); err != nil {
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

		if len(e.headerKeys) != len(e.header) {
			return errors.New("exportx: length of headerKeys is not equal to length of header")
		}

		// set data
		for i := 0; i < len(arr); i++ {
			rowNum := i + 2 // 从第二行开始写数据(第一行写header)
			for j := 0; j < len(e.headerKeys); j++ {
				coordinate, _ := excelize.CoordinatesToCellName(j+1, rowNum)
				if err = e.file.SetCellValue(e.sheetName, coordinate, arr[i][e.headerKeys[j]]); err != nil {
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

				coordinate, _ := excelize.CoordinatesToCellName(i+1, 1)
				if err = e.file.SetCellValue(e.sheetName, coordinate, tag); err != nil {
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

				coordinate, _ := excelize.CoordinatesToCellName(j+1, rowNum)
				if err = e.file.SetCellValue(e.sheetName, coordinate, ev); err != nil {
					return errors.Wrap(err, "failed to set data value")
				}
			}
		}
	default:
		return ErrInvalidDataFormat
	}

	return
}

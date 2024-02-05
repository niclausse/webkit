package exportx

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
	"log"
	"os"
)

var (
	ErrInvalidHeader     = errors.New("exportx: invalid header")
	ErrInvalidDataFormat = errors.New("exportx: invalid data format")
)

type CSV struct {
	*config
}

func NewCSV() Exporter {
	return &CSV{
		config: &config{},
	}
}

func (c *CSV) Export(data interface{}, opts ...Option) ([]byte, error) {
	for _, opt := range opts {
		opt(c.config)
	}

	return c.genCSV(data)
}

func (c *CSV) ExportLocal(data interface{}, opts ...Option) (string, error) {
	for _, opt := range opts {
		opt(c.config)
	}

	if err := c.initSaveProfile(CsvExt); err != nil {
		return "", err
	}

	var (
		err     error
		content []byte
		file    *os.File
	)

	content, err = c.genCSV(data)
	if err != nil {
		return "", err
	}

	file, err = os.Create(c.savePath)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}(file)

	if _, err = file.Write(content); err != nil {
		return "", err
	}

	return c.savePath, nil
}

func (c *CSV) writeStructSlice(data interface{}) error {
	file, err := os.OpenFile(c.savePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "fail to open file %s", c.savePath)
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}(file)

	b, err := csvutil.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "fail to marshal data")
	}

	buf := bufio.NewWriter(file)

	_, err = buf.Write(b)
	if err != nil {
		return errors.Wrap(err, "fail to write data")
	}

	err = buf.Flush()
	if err != nil {
		return errors.Wrap(err, "fail to flush data")
	}

	return nil
}

func (c *CSV) genCSV(data interface{}) ([]byte, error) {
	var (
		err error
		w   bytes.Buffer
	)

	switch v := data.(type) {
	case [][]string:
		if len(c.header) == 0 {
			return nil, ErrInvalidHeader
		}

		if _, err = w.WriteString("\xEF\xBB\xBF"); err != nil { // write UTF-8 BOM
			return nil, err
		}

		writer := csv.NewWriter(&w) // csv writer

		if err = writer.Write(c.header); err != nil {
			return nil, errors.Wrap(err, "fail to write csv header")
		}

		if err = writer.WriteAll(v); err != nil {
			return nil, errors.Wrap(err, "fail to write csv raw data")
		}

		return w.Bytes(), nil
	default:
		var content []byte
		content, err = csvutil.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(err, "fail to marshal data")
		}

		buf := bufio.NewWriter(&w)

		_, err = buf.Write(content)
		if err != nil {
			return nil, errors.Wrap(err, "fail to write data")
		}

		if err = buf.Flush(); err != nil {
			return nil, errors.Wrap(err, "fail to flush data")
		}

		return w.Bytes(), nil
	}
}

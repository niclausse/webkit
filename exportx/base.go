package exportx

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Exporter interface {
	Export(data interface{}, opts ...Option) ([]byte, error)
	ExportLocal(data interface{}, opts ...Option) (string, error)
	Close() error
}

type config struct {
	header           []string
	expectedFilename string // expected filename to save
	savedFilename    string // real filename to save
	baseDir          string // basic dir to save
	dir              string // unique path to save each file
	savePath         string // file's full path,  usually: dir + savedFilename

}

type Option func(cfg *config)

func WithFilename(name string) Option {
	return func(cfg *config) {
		cfg.expectedFilename = name
	}
}

func WithHeader(header []string) Option {
	return func(cfg *config) {
		cfg.header = header
	}
}

func WithBaseDir(dir string) Option {
	return func(cfg *config) {
		cfg.baseDir = dir
	}
}

type (
	Extend string
	Type   string
)

func (e Extend) String() string {
	return string(e)
}

const (
	ExcelStorePath = "./export_files/excel/"
	CsvStorePath   = "./export_files/csv/"

	ExcelExt Extend = ".xlsx"
	CsvExt   Extend = ".csv"
	ZipExt   Extend = ".zip"

	ExcelType Type = "excel"
	CsvType   Type = "csv"

	ExcelLimit = 30000
	CsvLimit   = 50000
)

func generateUniqueDirName() string {
	return fmt.Sprintf("%s-%d/", uuid.NewString(), time.Now().UnixNano())
}

func generateUniqueFilename(index int, name string) string {
	return fmt.Sprintf("%d.%s-%s", index, time.Now().Format(time.RFC3339), name)
}

func appendFileExt(name string, ext Extend) string {
	return name + string(ext)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func (c *config) initSaveProfile(ext Extend) error {
	if c.expectedFilename == "" {
		c.savedFilename = generateUniqueFilename(1, appendFileExt(uuid.NewString(), ext))
	} else if strings.HasSuffix(c.expectedFilename, ext.String()) {
		c.savedFilename = generateUniqueFilename(1, c.expectedFilename)
	} else {
		c.savedFilename = generateUniqueFilename(1, appendFileExt(c.expectedFilename, ext))
	}

	if c.baseDir == "" {
		if ext == CsvExt {
			c.baseDir = CsvStorePath
		} else {
			c.baseDir = ExcelStorePath
		}
	}

	c.dir = filepath.Join(c.baseDir, generateUniqueDirName())
	if !isExist(c.dir) {
		oldMask := syscall.Umask(0)
		if err := os.MkdirAll(c.dir, os.ModePerm); err != nil {
			return err
		}
		syscall.Umask(oldMask)
	}

	c.savePath = filepath.Join(c.dir, c.savedFilename)
	return nil
}

package storage

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	defaultSavePath = "./data.csv"
)

type FileStorage struct {
	SavePath  string
	dataCache map[string]*Record
}

func NewFileStorage() Storage {
	cache := make(map[string]*Record)
	return &FileStorage{
		dataCache: cache,
	}
}

func (fs *FileStorage) Load() error {
	return nil
}

func (fs *FileStorage) Data() []*Record {
	res := make([]*Record, 0, len(fs.dataCache))
	for _, record := range fs.dataCache {
		res = append(res, record)
	}
	return res
}

func (fs *FileStorage) Save(record *Record) error {
	fs.dataCache[record.ID] = record

	rows := ""
	for _, record := range fs.dataCache {
		rows += record.String() + "\n"
	}

	if fs.SavePath == "" {
		fs.SavePath = defaultSavePath
	}

	file, err := os.OpenFile(fs.SavePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "open file %s faild", fs.SavePath)
	}
	defer file.Close()

	fmt.Fprintln(file, rows)

	return nil
}

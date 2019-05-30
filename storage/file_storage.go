package storage

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	defaultSavePath = "./data/data.csv"
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
	log.Println("[START] load data")

	fs.dataCache = map[string]*Record{}
	if fs.SavePath == "" {
		fs.SavePath = defaultSavePath
	}

	file, err := os.Open(fs.SavePath)
	if err != nil {
		return errors.Wrapf(err, "open file %s faild", fs.SavePath)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		record, err := NewFromCSVLine(line)
		if err != nil {
			return errors.Wrap(err, "error new from csv line")
		}
		fs.dataCache[record.ID] = record
	}
	log.Println("[DONE] load data")
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
	return fs.save()
}

func (fs *FileStorage) Delete(id string) error {
	delete(fs.dataCache, id)
	return fs.save()
}

func (fs *FileStorage) save() error {
	if fs.SavePath == "" {
		fs.SavePath = defaultSavePath
	}

	file, err := os.OpenFile(fs.SavePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Wrapf(err, "open file %s faild", fs.SavePath)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, record := range fs.dataCache {
		writer.Write(record.Strings())
	}
	writer.Flush()

	return nil
}

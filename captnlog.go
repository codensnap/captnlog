package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	bolt "go.etcd.io/bbolt"
)

const (
	filename = ".captainlog.bdb"
)

var (
	configBkt      = []byte("config")
	entryBkt       = []byte("entries")
	defaultBkt     = []byte("default")
	categorizedBkt = []byte("categorized")

	initBuckets = [][]byte{
		configBkt,
		entryBkt,
		categorizedBkt,
	}
)

var (
	ErrNoBucket = errors.New("bucket does not exist")
)

func defaultLocation() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory: %w", err)
	}
	filePath := path.Join(dir, filename)
	return filePath, nil
}

type CaptnLog struct {
	bdb *bolt.DB
	lgr hclog.Logger
}

func New() (*CaptnLog, error) {
	filePath, err := defaultLocation()
	if err != nil {
		return nil, err
	}
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to open bdb: %w", err)
	}
	lgr := hclog.Default()
	lgr.SetLevel(hclog.Debug)
	cl := &CaptnLog{
		bdb: db,
		lgr: lgr,
	}
	if err = cl.init(); err != nil {
		db.Close()
		return nil, err
	}
	return cl, nil
}

func (cl *CaptnLog) init() error {
	cl.bdb.Update(func(tx *bolt.Tx) error {
		for _, b := range initBuckets {
			tx.CreateBucketIfNotExists(b)
		}
		return nil
	})
	return nil
}

func (cl *CaptnLog) CommandFactory(cmdId CmdID, category string) func() (cli.Command, error) {
	var clc CaptnLogCommand
	switch cmdId {
	case WriteCmd:
		clc = WriteCommand
	case ReadCmd:
		clc = ReadCommand
	case ReadAllCmd:
		clc = ReadAllCommand
	default:
		panic("No")
	}
	clc.captnLog = cl
	clc.category = category
	return func() (cli.Command, error) {
		return &clc, nil
	}

}

func (cl *CaptnLog) WriteEntry(category, entry string) error {
	if entry == "" {
		return nil
	}
	log := &Log{
		Timestamp: time.Now(),
		Category:  category,
		Entry:     entry,
	}
	cl.lgr.Debug("writing entry",
		"category", category,
		"entry", entry)
	return cl.bdb.Update(func(tx *bolt.Tx) error {
		entryBkt := tx.Bucket(entryBkt)
		if entryBkt == nil {
			return ErrNoBucket
		}
		bkt := entryBkt.Bucket([]byte(category))
		if bkt == nil {
			newBkt, err := entryBkt.CreateBucket([]byte(category))
			if err != nil {
				return err
			}
			bkt = newBkt
		}
		b, err := log.Encode()
		if err != nil {
			return err
		}
		err = bkt.Put(log.Key(), b)
		if err != nil {
			cl.lgr.Error(err.Error())
		}
		return err
	})
}

func (cl *CaptnLog) ReadEntries(category string) ([]*Log, error) {
	logs := []*Log{}
	cl.lgr.Debug("reading entries in " + category)
	err := cl.bdb.View(func(tx *bolt.Tx) error {
		entryBkt := tx.Bucket(entryBkt)
		if entryBkt == nil {
			return ErrNoBucket
		}
		bkt := entryBkt.Bucket([]byte(category))
		if bkt == nil {
			return ErrNoBucket
		}
		bkt.ForEach(func(k, v []byte) error {
			log, err := DecodeLog(v)
			if err != nil {
				cl.lgr.Error(err.Error())
				return err
			}
			logs = append(logs, log)
			return nil
		})
		return nil
	})
	return logs, err
}

func (cl *CaptnLog) ReadAllEntries() ([]*Log, error) {
	cl.lgr.Debug("reading all")
	logs := []*Log{}
	err := cl.bdb.View(func(tx *bolt.Tx) error {
		entryBkt := tx.Bucket(entryBkt)
		if entryBkt == nil {
			return ErrNoBucket
		}
		entryBkt.ForEach(func(k, v []byte) error {
			if v != nil {
				return nil
			}
			bkt := entryBkt.Bucket(k)
			return bkt.ForEach(func(k, v []byte) error {
				log, err := DecodeLog(v)
				if err != nil {
					return err
				}
				logs = append(logs, log)
				return nil
			})
		})
		return nil
	})
	return logs, err
}

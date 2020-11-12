package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Log struct {
	Timestamp time.Time
	Category  string
	Entry     string
}

func (l *Log) Key() []byte {
	return []byte(l.Timestamp.Format(time.RFC3339))
}

func (l *Log) Encode() ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(*l); err != nil {
		return nil, fmt.Errorf("failed to encode log: %w", err)
	}
	return b.Bytes(), nil
}

func DecodeLog(b []byte) (*Log, error) {
	reader := bytes.NewReader(b)
	decoder := gob.NewDecoder(reader)
	var log Log
	if err := decoder.Decode(&log); err != nil {
		return nil, fmt.Errorf("failed to decode log: %w", err)
	}
	return &log, nil
}

type Logs []*Log

func (l Logs) Len() int {
	return len(l)
}

func (l Logs) Less(i, j int) bool {
	return l[i].Timestamp.Before(l[j].Timestamp)
}

func (l Logs) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

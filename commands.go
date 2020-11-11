package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

var (
	ErrInvalidCommand = errors.New("invalid command")
)

type CmdID uint

const (
	WriteCmd CmdID = iota
	ReadCmd
)

type CaptnLogCommand struct {
	captnLog *CaptnLog
	category string
	help     string
	synopsis string
	cmdId    CmdID
}

func (clc *CaptnLogCommand) Help() string {
	return clc.help
}

func (clc *CaptnLogCommand) Synopsis() string {
	return clc.synopsis
}

func (clc *CaptnLogCommand) Run(args []string) int {
	var err error
	switch clc.cmdId {
	case WriteCmd:
		err = clc.captnLog.WriteEntry(clc.category, args[0])
	case ReadCmd:
		logs, err := clc.captnLog.ReadEntries(clc.category)
		if err != nil {
			return 1
		}
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)
		fmt.Fprintln(tw, "timestamp\tentry\t")
		for _, log := range logs {
			fmt.Fprintf(tw, "%s\t%s\t\n", log.Timestamp.Format(time.Stamp), log.Entry)
		}
		tw.Flush()
	}
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

var (
	ReadCommand = CaptnLogCommand{
		help:     "TODO",
		synopsis: "read entries in your captain's log",
		cmdId:    ReadCmd,
	}
	WriteCommand = CaptnLogCommand{
		help:     "TODO",
		synopsis: "write an entry to your captain's log",
		cmdId:    WriteCmd,
	}
)

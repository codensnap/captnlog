package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/mitchellh/cli"
)

var (
	ErrInvalidCommand = errors.New("invalid command")

	ui = cli.BasicUi{}
)

type CmdID uint

const (
	WriteCmd CmdID = iota
	ReadCmd
	ReadAllCmd
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
			ui.Error(err.Error())
			return 1
		}
		printLogs(logs)
	case ReadAllCmd:
		logs, err := clc.captnLog.ReadAllEntries()
		if err != nil {
			ui.Error(err.Error())
			return 1
		}
		printLogs(logs)
	}
	if err != nil {
		ui.Error(err.Error())
		return 1
	}
	return 0
}

func printLogs(logs Logs) {
	sort.Sort(logs)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "timestamp\tentry\tcategory\t")
	for _, log := range logs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t\n", log.Timestamp.Format(time.Stamp), log.Entry, log.Category)
	}
	tw.Flush()
}

var (
	ReadCommand = CaptnLogCommand{
		help:     "TODO",
		synopsis: "read entries from your captain's log",
		cmdId:    ReadCmd,
	}
	WriteCommand = CaptnLogCommand{
		help:     "TODO",
		synopsis: "write an entry to your captain's log",
		cmdId:    WriteCmd,
	}
	ReadAllCommand = CaptnLogCommand{
		help:     "TODO",
		synopsis: "read all entries in all categories",
		cmdId:    ReadAllCmd,
	}
)

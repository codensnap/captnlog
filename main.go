package main

import (
	"flag"

	"github.com/mitchellh/cli"
)

var (
	category = flag.String("c", "default", "select a category")
)

func main() {
	flag.Parse()
	cl, err := New()
	if err != nil {
		panic(err)
	}
	c := cli.NewCLI("captnlog", "0.0.1")
	c.Args = flag.Args()
	c.Commands = map[string]cli.CommandFactory{
		"write": cl.CommandFactory(WriteCmd, *category),
		"read":  cl.CommandFactory(ReadCmd, *category),
	}
	if _, ok := c.Commands[c.Subcommand()]; !ok {
	}
	c.Run()
}

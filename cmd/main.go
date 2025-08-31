package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const supportedCommands = "'add', 'list'"
const defaultStoreFile = "./store.json"

func main() {
	if len(os.Args) < 2 {
		exitWithErrorMessage(fmt.Sprintf("arguments. Needed. Supported commands %s", supportedCommands))
	}

	command := os.Args[1]

	switch command {
	case "add":
		var name, nws string
		var amount int
		addFlags := flag.NewFlagSet("add", flag.ExitOnError)
		addFlags.StringVar(&name, "name", "", "Name of the expense")
		addFlags.StringVar(&nws, "nws", "", "Category of the expense: needs/wants/savings")
		addFlags.IntVar(&amount, "amount", -1, "Sum spent")
		addFlags.Parse(os.Args[2:])
		nws = strings.ToLower(nws)
		if err := validateArgs(name, nws, amount); err != nil {
			exitWithErrorMessage(err.Error())
		}
	case "list":
		if len(os.Args[2:]) > 0 {
			exitWithErrorMessage("list: no additional arguments supported")
		}
	default:
		exitWithErrorMessage(fmt.Sprintf("wrong command %s. Supported commands %s", command, supportedCommands))
	}
}

func validateArgs(name, nws string, amount int) error {
	if name == "" {
		return errors.New("--name cannot be empty")
	}

	if nws == "" {
		return errors.New("--nws cannot be empty")
	}

	if nws != "wants" && nws != "needs" && nws != "savings" {
		return errors.New("supported values for --nws: needs/wants/savings")
	}

	if amount <= 0 {
		return errors.New("--amount is required and needs to be positive")
	}

	return nil
}

func exitWithErrorMessage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

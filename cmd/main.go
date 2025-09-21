package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/efistokl/nws-ledger/operations"
)

const supportedCommands = "'add', 'list', 'summary'"
const DefaultStoreFile = "./store.json"

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("arguments needed. Supported commands %s", supportedCommands)
	}

	database, err := os.OpenFile(DefaultStoreFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("error opening storage file %s: %v", DefaultStoreFile, err)
	}
	defer database.Close()

	store, err := operations.NewJSONStorage(database)

	if err != nil {
		log.Fatalf("failed creating storage from file %s: %v", DefaultStoreFile, err)
	}

	command := os.Args[1]

	switch command {
	case "add":
		var name, nws, domain string
		var amount int
		addFlags := flag.NewFlagSet("add", flag.ExitOnError)
		addFlags.StringVar(&name, "name", "", "Name of the expense")
		addFlags.StringVar(&nws, "nws", "", "Category of the expense: needs/wants/savings")
		addFlags.StringVar(&domain, "domain", "", "Domain category of the expense: groceries/shopping/rent/... (whatever)")
		addFlags.IntVar(&amount, "amount", -1, "Sum spent")
		addFlags.Parse(os.Args[2:])
		nws = strings.ToLower(nws)
		if err := validateArgs(name, nws, amount); err != nil {
			log.Fatal(err.Error())
		}

		store.Add(operations.Expense{Amount: amount, Name: name, NWS: nws, Domain: domain})
	case "list":
		if len(os.Args[2:]) > 0 {
			log.Fatal("list: no additional arguments supported")
		}

		fmt.Print(operations.FormatCSVList(store))
	case "summary":
		if len(os.Args[2:]) > 0 {
			log.Fatal("summary: no additional arguments supported")
		}

		fmt.Print(operations.FormatCSVSummaryByNWS(store))
	default:
		log.Fatalf("wrong command %s. Supported commands %s", command, supportedCommands)
	}
}

func validateArgs(name, nws string, amount int) error {
	if name == "" {
		return errors.New("--name cannot be empty")
	}

	if nws == "" {
		return errors.New("--nws cannot be empty")
	}

	if err := operations.ValidateNWS(nws); err != nil {
		return fmt.Errorf("--nws: %w", err)
	}

	if amount <= 0 {
		return errors.New("--amount is required and needs to be positive")
	}

	return nil
}

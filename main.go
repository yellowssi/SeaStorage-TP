package main

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/jessevdk/go-flags"
	"gitlab.com/SeaStorage/SeaStorage-TP/handler"
	"os"
	"syscall"
)

const (
	FamilyName    string = "SeaStorage"
	FamilyVersion string = "1.0.0"
)

type Opts struct {
	Version bool   `short:"V" long:"version" description:"Display version"`
	Verbose []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	Connect string `short:"C" long:"connect" description:"Validator component endpoint to connect to" default:"tcp://localhost:4004"`
}

func main() {
	var opts Opts

	logger := logging.Get()

	parser := flags.NewParser(&opts, flags.Default)
	remaining, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			logger.Error("Failed to parse args: %v", err)
			os.Exit(2)
		}
	}

	if len(remaining) > 0 {
		fmt.Printf("Error: Unrecognized arguments passed: %v\n", remaining)
		os.Exit(2)
	}

	if opts.Version {
		fmt.Println("SeaStorage Transaction Processor")
		fmt.Println("Version: " + FamilyVersion)
		os.Exit(0)
	}

	endpoint := opts.Connect

	switch len(opts.Verbose) {
	case 2:
		logger.SetLevel(logging.DEBUG)
	case 1:
		logger.SetLevel(logging.INFO)
	default:
		logger.SetLevel(logging.WARN)
	}

	logger.Debugf("command line arguments: %v\n", os.Args)
	logger.Debugf("verbose = %v\n", len(opts.Verbose))
	logger.Debugf("endpoint = %v\n", endpoint)

	hd := handler.NewSeaStorageHandler(FamilyName, []string{FamilyVersion})
	proc := processor.NewTransactionProcessor(endpoint)
	proc.AddHandler(hd)
	proc.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	err = proc.Start()
	if err != nil {
		logger.Error("Processor stopped: ", err)
	}
}

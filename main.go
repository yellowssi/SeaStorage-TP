package main

import (
	"encoding/gob"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/jessevdk/go-flags"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/seaStorageHandler"
	"gitlab.com/SeaStorage/SeaStorage-Hyperledger/pkg/storage"
	"os"
	"syscall"
)

func init() {
	gob.Register(&storage.File{})
	gob.Register(&storage.Directory{})
}

type Opts struct {
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

	handler := seaStorageHandler.NewSeaStorageHandler([]string{"1.0"})
	proc := processor.NewTransactionProcessor(endpoint)
	proc.AddHandler(handler)
	proc.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	err = proc.Start()
	if err != nil {
		logger.Error("Processor stopped: ", err)
	}
}

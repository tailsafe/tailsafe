package main

import (
	"flag"
	"github.com/tailsafe/tailsafe/internal/tailsafe"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	tailsafeInterface "github.com/tailsafe/tailsafe/pkg/tailsafe"
)

var workflow bool
var use string
var env string
var verbose bool
var data string

func main() {

	flag.BoolVar(&verbose, "v", false, "Show more information")
	flag.BoolVar(&workflow, "workflow", false, "Show the processing workflow")
	flag.StringVar(&data, "data", "", "Define the path name of the mock data")
	flag.StringVar(&use, "use", "", "Set use path name")
	flag.StringVar(&env, "env", "", "environment arguments")
	flag.Parse()

	// Register modules
	modules.Register("Utils", modules.GetUtilsModule())
	modules.Register("Events", modules.GetEventsModule())
	modules.Register("Logger", modules.GetLoggerModule())
	modules.Register("AsyncQueue", modules.GetAsyncQueue())

	// set user flags
	if workflow {
		modules.GetLoggerModule().AddNamespace(tailsafeInterface.NAMESPACE_WORKFLOW)
	}
	if verbose {
		modules.GetLoggerModule().SetVerbose(verbose)
	}

	// Create tailsafe-cli instance and run !
	tailsafe.
		New().
		SetPath(use).
		SetEnv(env).
		SetDataPath(data).
		Run()
}

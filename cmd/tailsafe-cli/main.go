package main

import (
	"flag"
	"github.com/tailsafe/tailsafe/internal/tailsafe"
	"github.com/tailsafe/tailsafe/internal/tailsafe/modules"
	tailsafeInterface "github.com/tailsafe/tailsafe/pkg/tailsafe"
)

var workflow bool
var path string
var env string
var verbose bool
var data string

func main() {

	flag.BoolVar(&verbose, "v", false, "Show more information")
	flag.BoolVar(&workflow, "workflow", false, "Show the workflow")
	flag.StringVar(&data, "data", "", "Set specific data path name")
	flag.StringVar(&path, "path", "", "Set config path name")
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
		SetPath(path).
		SetEnv(env).
		SetPathData(data).
		Run()
}

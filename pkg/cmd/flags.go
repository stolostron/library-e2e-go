package cmd

import (
	"flag"
)

type End2EndT struct {
	OptionsFile string
	ReportFile  string
	Owner       string
	UID         string
}

var End2End End2EndT

//InitFlags adds some default flags for the e2e test command-line
func InitFlags(flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.CommandLine
	}
	flagset.StringVar(&End2End.Owner, "owner", "",
		"Provide the prefix for created resources (e.g. -owner=\"xxxxx\").\n"+
			"If not present the owner defined in the options.yaml will be taken.\n"+
			"If not present the environment variable $UESR will be taken.")
	flagset.StringVar(&End2End.UID, "uid", "",
		"Provide a unique ID as postfix (e.g. -uid=\"xxxx\").\n"+
			"If not present 4 random chars will be generated")
	flagset.StringVar(&End2End.ReportFile, "report-file", "results",
		"Provide the path to where the junit results will be printed.")
	flagset.StringVar(&End2End.OptionsFile, "options", "",
		"Location of an \"options.yaml\" file to provide input for various tests")
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/op/go-logging"
)

// later, for non-AWS certs:
// echo | timeout 5 openssl s_client -servername basket.mozilla.org -connect basket.mozilla.org111:443 2>/dev/null | openssl x509 -noout -dates

var log = logging.MustGetLogger("heimdall")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func setupLogging(enableDebug bool) {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	if enableDebug {
		backend1Leveled.SetLevel(logging.DEBUG, "")
	} else {
		backend1Leveled.SetLevel(logging.INFO, "")
	}
	logging.SetBackend(backend1Leveled)
}

func main() {
	var regions arrayFlags
	flag.Var(&regions, "region", "Region to scan")
	debug := flag.Bool("debug", false, "Show debugging output")

	flag.Parse()

	setupLogging(*debug)

	if len(regions) == 0 {
		fmt.Println("Please specify a region")
		os.Exit(1)
	}
	var allRegions []*RegionTests
	for _, region := range regions {
		log.Infof("Checking region %s\n", region)
		if result, err := DoScan(region); err != nil {
			fmt.Println("FOO")
		} else {
			allRegions = append(allRegions, result)
		}
	}

	b, err := json.Marshal(allRegions)
	if err != nil {
		log.Error("error:", err)
	}
	os.Stdout.Write(b)
}

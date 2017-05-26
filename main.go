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
	`%{message}`,
	// fancy color output:
	//`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func setupLogging(logLevel logging.Level) {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logLevel, "")
	logging.SetBackend(backend1Leveled)
}

func main() {
	var regions arrayFlags
	flag.Var(&regions, "ec2-region", "EC2 Region to scan ELBs")
	hostsFile := flag.String("hosts", "", "Hosts file to process")
	warnDays := flag.Int("warn-days", 60, "Warn on certs with <= this value to expiration")
	logLevel := flag.String("log-level", "error", "Log level (debug, info, error)")
	dumpJson := flag.Bool("json", false, "Display query results as JSON")
	skipExpired := flag.Bool("skip-expired", false, "Skip certs that have already expired")

	flag.Parse()

	logLevels := map[string]logging.Level{
		"debug": logging.DEBUG,
		"info":  logging.INFO,
		"error": logging.ERROR,
	}

	if val, ok := logLevels[*logLevel]; ok {
		setupLogging(val)
	} else {
		fmt.Println("Invalid log-level")
		os.Exit(1)
	}

	// todo: move this to scan.go
	var allRegions []*RegionTests
	for _, region := range regions {
		log.Infof("Checking region %s\n", region)
		if result, err := processRegionELBs(region); err != nil {
			log.Error(err)
		} else {
			allRegions = append(allRegions, result)
		}
	}
	if len(regions) > 0 {
		showManagedExpirations(allRegions, *warnDays, *skipExpired)
	}

	// this is lame
	if *hostsFile != "" {
		checkUnmanagedExpirations(*hostsFile, *warnDays, *skipExpired)
	}

	// move this to it's own fn, or just drop it
	if *dumpJson {
		b, err := json.Marshal(allRegions)
		if err != nil {
			log.Error("error:", err)
		}
		os.Stdout.Write(b)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/melroy89/angie_telegraf_plugin/plugins/inputs/angie_api"

	"github.com/influxdata/telegraf/plugins/common/shim"
)

var pollInterval = flag.Duration("poll_interval", 10*time.Second, "How often to send metrics")
var pollIntervalEnabled = flag.Bool("poll_interval_enabled", false, "Enable polling interval")
var configFile = flag.String("config", "", "path to the config file for this plugin")
var err error

func main() {
	// parse command line options
	flag.Parse()
	if !*pollIntervalEnabled {
		*pollInterval = shim.PollIntervalDisabled
	}

	// create the shim. This is what will run your plugins.
	s := shim.New()

	// Check for settings from a config toml file,
	err = s.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err loading input: %s\n", err)
		os.Exit(1)
	}

	// run the plugin until stdin closes or we receive a termination signal
	if err := s.Run(*pollInterval); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}

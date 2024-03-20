package main

import (
	"flag"
	"log"

	"github.com/h0tak88r/subMonit88r/config"
	"github.com/h0tak88r/subMonit88r/runner"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	runner.PrintLogo()

	// Define flag variables
	var domainList string
	var webhook string
	var monitor bool

	// Parse flags
	flag.StringVar(&domainList, "l", "", "Specify a file containing a list of domains")
	flag.StringVar(&webhook, "wh", "", "Specify the Discord webhook URL")
	flag.BoolVar(&monitor, "m", false, "Enable subdomain monitoring")
	flag.Parse()

	// Set the configuration values
	config.SetConfig(domainList, webhook, monitor)

	// Validate flags
	if err := config.ValidateFlags(); err != nil {
		log.Fatal(err)
	}

	runner.Run()
}

package config

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Exported Constants
const (
	DbFile          = "subdomains_database.db"
	ResultsFileName = "subfalconResults.txt"
	MonitorInterval = 5 * time.Hour
)

// Config holds the configuration options.
type Config struct {
	DomainList   string
	Webhook      string
	Monitor      bool
	SingleDomain string // New field to store the single domain
}

var (
	cfg        Config
	configLock sync.Mutex
)

var (
	ErrMissingDomainListFlag = errors.New("missing domain list flag or single domain")
)

// PrintLogo prints the logo.
func PrintLogo() {
	fmt.Println(`
	┏┓  ┓ ┏  ┓       
	┗┓┓┏┣┓╋┏┓┃┏┏┓┏┓  ·˚ * 🔭 ⋆ .☆ 
	┗┛┗┻┗┛┛┗┻┗┗┗┛┛┗ 
			By: @h0tak88r @zomasec
	`)
}

// SetConfig sets the configuration options and returns a pointer to the updated Config.
func SetConfig(domainList, webhook string, monitor bool, singleDomain string) *Config {
	configLock.Lock()
	defer configLock.Unlock()

	cfg = Config{
		DomainList:   domainList,
		Webhook:      webhook,
		Monitor:      monitor,
		SingleDomain: singleDomain, // Set the single domain in the config
	}

	return &cfg
}

// GetConfig returns a pointer to the current configuration options.
func GetConfig() *Config {
	configLock.Lock()
	defer configLock.Unlock()

	return &cfg
}

// ValidateFlags validates the required flags.
func ValidateFlags() error {
	// Ensure either a domain list or a single domain is provided
	if cfg.DomainList == "" && cfg.SingleDomain == "" {
		return ErrMissingDomainListFlag
	}
	return nil
}

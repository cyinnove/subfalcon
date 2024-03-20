package config

import (
	"errors"
	"sync"
	"time"
)

// Exported Constants
const (
	DbFile          = "subdomains_database.db"
	ResultsFileName = "subMonit88rResults.txt"
	MonitorInterval = 5 * time.Hour
)

// Config holds the configuration options.
type Config struct {
	DomainList string
	Webhook    string
	Monitor    bool
}

var (
	cfg        Config
	configLock sync.Mutex
)

var ErrMissingDomainListFlag = errors.New("missing domain list flag")

// SetConfig sets the configuration options and returns a pointer to the updated Config.
func SetConfig(domainList string, webhook string, monitor bool) *Config {
	configLock.Lock()
	defer configLock.Unlock()

	cfg = Config{
		DomainList: domainList,
		Webhook:    webhook,
		Monitor:    monitor,
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
	if cfg.DomainList == "" {
		return errors.New("missing domain list flag")
	}
	return nil
}

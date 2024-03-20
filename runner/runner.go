package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/h0tak88r/subMonit88r/config"
	"github.com/h0tak88r/subMonit88r/pkg/db"
	"github.com/h0tak88r/subMonit88r/pkg/sub88r"
)

var cfg = config.GetConfig()

// Run is the main entry point for the runner package.
func Run() {
	// Start running
	go db.InitDB(config.DbFile)

	if cfg.Monitor {
		for {
			fmt.Println("[+] Monitoring subdomains in domains.txt file......")
			PassiveSubdomainEnumeration()
			time.Sleep(config.MonitorInterval)
		}
	} else {
		PassiveSubdomainEnumeration()
	}
}

func PrintLogo() {
	fmt.Println(`
		┏┓  ┓ ┳┳┓    • ┏┓┏┓  
		┗┓┓┏┣┓┃┃┃┏┓┏┓┓╋┣┫┣┫┏┓
		┗┛┗┻┗┛┛ ┗┗┛┛┗┗┗┗┛┗┛┛ 
				by sallam(h0tak88r)
	`)
}

func PassiveSubdomainEnumeration() {
	domains, err := readDomainsFromFile(cfg.DomainList)
	if err != nil {
		log.Fatal(err)
	}

	uniqueSubdomains := make(map[string]struct{})

	for _, domain := range domains {
		subdomains := fetchSubdomainsFromSources(domain)
		for _, subdomain := range subdomains {
			uniqueSubdomains[subdomain] = struct{}{}
		}
	}

	var allSubdomains []string
	for subdomain := range uniqueSubdomains {
		allSubdomains = append(allSubdomains, subdomain)
	}

	oldSubdomains := db.Getsubdomains(config.DbFile)
	newSubdomains := difference(allSubdomains, oldSubdomains)

	writeSubdomainsToFile(config.ResultsFileName, allSubdomains)

	fmt.Println("[+] Subdomains Enumeration completed, Results are saved in subMonit88rResults.txt.")

	if len(newSubdomains) > 0 {
		fmt.Printf("[+] %d new subdomains discovered:\n", len(newSubdomains))
		db.AddSubdmomains(newSubdomains, config.DbFile)

		for _, subdomain := range newSubdomains {
			fmt.Println(subdomain)
		}
		// Notify user via Discord webhook if provided
		if cfg.Webhook != "" {
			go sendDiscordNotification(newSubdomains)
		}
	}
}

func sendDiscordNotification(subdomains []string) {
	message := fmt.Sprintf("[Subdomain Monitor] %d new subdomains discovered:\n", len(subdomains))
	for _, subdomain := range subdomains {
		message += subdomain + "\n"
	}

	if err := sendWebhookMessage(cfg.Webhook, message); err != nil {
		fmt.Printf("[!] Error sending Discord notification: %v\n", err)
	}
}

func sendWebhookMessage(webhook, message string) error {
	payload := map[string]string{"content": message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("non-OK status code: %d", resp.StatusCode)
	}

	return nil
}

// Fetch subdomains from different sources using subber package
func fetchSubdomainsFromSources(domain string) []string {
	var wg sync.WaitGroup

	// Create a Subber instance
	subber := &sub88r.Subber{
		Domain:  domain,
		Results: &sub88r.Results{},
	}

	// Define a function to fetch subdomains from a source
	fetch := func(fetchFunc func() error, sourceName, domain string) {
		defer wg.Done()
		fmt.Printf("[+] Getting Subdomains from %s for %s...\n", sourceName, domain)
		if err := fetchFunc(); err != nil {
			log.Fatalf("Error while getting subdomains from %s for %s: %v\n", sourceName, domain, err)
		}
	}

	// Fetch subdomains from each source concurrently
	wg.Add(5)
	go fetch(subber.Anubis, "Anubis jdlc.me", domain)
	go fetch(subber.UrlScan, "UrlScan", domain)
	go fetch(subber.CrtSh, "CrtSh", domain)
	go fetch(subber.HackerTarget, "HackerTarget", domain)
	go fetch(subber.Otx, "Otx", domain)

	wg.Wait()

	// Return the combined subdomains
	return subber.GetAllSubdomains()
}

func difference(set1, set2 []string) []string {
	var diff []string
	set := make(map[string]struct{})

	for _, s := range set2 {
		set[s] = struct{}{}
	}

	for _, s := range set1 {
		if _, ok := set[s]; !ok {
			diff = append(diff, s)
		}
	}

	return diff
}

func writeSubdomainsToFile(filename string, subdomains []string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	uniqueSubdomains := make(map[string]struct{})

	for _, subdomain := range subdomains {
		if _, ok := uniqueSubdomains[subdomain]; !ok {
			file.WriteString(subdomain + "\n")
			uniqueSubdomains[subdomain] = struct{}{}
		}
	}
}

func readDomainsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	data := []string{}

	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

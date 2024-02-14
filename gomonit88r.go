package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile          = "subdomains_database.db"
	resultsFileName = "subMonit88rResults.txt"
)

type config struct {
	DomainList string
	Webhook    string
	Monitor    bool
}

var cfg config

func init() {
	rootCmd := &cobra.Command{
		Use:   "goMonit88r",
		Short: "A tool for subdomain enumeration with monitoring and Discord notification support.",
		Run:   run,
	}
	rootCmd.Flags().StringVarP(&cfg.DomainList, "domain-list", "l", "", "Specify a file containing a list of domains")
	rootCmd.Flags().StringVarP(&cfg.Webhook, "webhook", "w", "", "Specify the Discord webhook URL")
	rootCmd.Flags().BoolVarP(&cfg.Monitor, "monitor", "m", false, "Enable subdomain monitoring")

	// Add validation for required flags
	err := rootCmd.MarkFlagRequired("domain-list")
	if err != nil {
		log.Fatal(err)
	}

	// Add custom validation for flags
	cobra.OnInitialize(validateFlags)

	rootCmdInstance = rootCmd
}

func validateFlags() {
	if cfg.DomainList == "" {
		fmt.Println("Error: Missing required flag --domain-list")
		fmt.Println("Use 'goMonit88r --help' for usage information.")
		os.Exit(1)
	}
}

var rootCmdInstance *cobra.Command
var notificationBuffer = make(chan string)
var monitorInterval = 10 * time.Hour

func printLogo() {
	fmt.Println(`
        ┳┳┓    • ┏┓┏┓  
    ┏┓┏┓┃┃┃┏┓┏┓┓╋┣┫┣┫┏┓
    ┗┫┗┛┛ ┗┗┛┛┗┗┗┗┛┗┛┛ 
     ┛    by sallam(h0tak88r)             
	`)

}
func main() {
	printLogo()
	err := rootCmdInstance.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	go printLogo()
	go initDatabase()

	if cfg.Monitor {
		monitorSubdomains()
	} else {
		subdomainEnumeration()
	}
}

func monitorSubdomains() {
	for {
		subdomainEnumeration()
		time.Sleep(monitorInterval)
	}
}

func subdomainEnumeration() {
	startTime := time.Now()
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

	oldSubdomains := getSubdomainsFromDB()
	newSubdomains := difference(allSubdomains, oldSubdomains)

	writeSubdomainsToFile(resultsFileName, allSubdomains)

	elapsedTime := time.Since(startTime)
	fmt.Printf("[+] Subdomains Enumeration completed in %s, Results are saved in subMonit88rResults.txt.\n", elapsedTime)

	if len(newSubdomains) > 0 {
		fmt.Printf("[+] %d new subdomains discovered:\n", len(newSubdomains))
		addSubdomainsToDB(newSubdomains)

		for _, subdomain := range newSubdomains {
			fmt.Println(subdomain)
		}
		// Notify user via Discord webhook if provided
		if cfg.Webhook != "" {
			go sendDiscordNotification(newSubdomains)
		}
	}
}

func initDatabase() {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS subdomains (id INTEGER PRIMARY KEY, subdomain TEXT)")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("[+] Database initialized successfully.")
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK status code: %d", resp.StatusCode)
	}

	return nil
}

func addSubdomainsToDB(subdomains []string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO subdomains (subdomain) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, subdomain := range subdomains {
		_, err = stmt.Exec(subdomain)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] Added %d new subdomains to the database.\n", len(subdomains))
}

func getSubdomainsFromDB() []string {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT subdomain FROM subdomains")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var subdomains []string

	for rows.Next() {
		var subdomain string
		err := rows.Scan(&subdomain)
		if err != nil {
			log.Fatal(err)
		}
		subdomains = append(subdomains, subdomain)
	}

	fmt.Println("[+] Retrieved subdomains from the database.")

	return subdomains
}

func fetchSubdomainsFromSources(domain string) []string {
	var subdomains []string

	// Fetch subdomains from various sources
	crtshSubdomains, _ := fetchSubdomainsFromCRTSH(domain)
	subdomains = append(subdomains, fetchSubdomainsFromAlienvault(domain)...)
	subdomains = append(subdomains, fetchSubdomainsFromAnubis(domain)...)
	subdomains = append(subdomains, fetchSubdomainsFromHackerTarget(domain)...)
	subdomains = append(subdomains, fetchSubdomainsFromRapiddns(domain)...)
	subdomains = append(subdomains, fetchSubdomainsFromUrlscan(domain)...)
	subdomains = append(subdomains, crtshSubdomains...)

	return subdomains
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

	var domains []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

// Fetch subdomains from various sources
func fetchSubdomainsFromCRTSH(domain string) ([]string, []string) {
	var subdomains []string
	var wildcardSubdomains []string

	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)
	fmt.Printf("[#] Fetching Subdomains from crt.sh for %s\n", domain)

	// Use gorequest to make HTTP request
	request := gorequest.New()
	_, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Printf("[!] Error fetching subdomains for domain %s from %s\n", domain, url)
		return subdomains, wildcardSubdomains
	}

	// Parse JSON response
	var entries []map[string]interface{}
	err := json.Unmarshal([]byte(body), &entries)

	if err != nil {
		fmt.Printf("[!] Error decoding JSON for domain %s from %s\n", domain, url)
		return subdomains, wildcardSubdomains
	}

	for _, entry := range entries {
		nameValue, ok := entry["name_value"].(string)
		if !ok {
			continue
		}

		if strings.Contains(nameValue, "\n") {
			subnameValues := strings.Split(nameValue, "\n")
			for _, subname := range subnameValues {
				subname = strings.TrimSpace(subname)
				if subname != "" {
					if strings.Contains(subname, "*") {
						wildcardSubdomains = append(wildcardSubdomains, subname)
					} else {
						subdomains = append(subdomains, subname)
					}
				}
			}
		}
	}

	return subdomains, wildcardSubdomains
}

func fetchSubdomainsFromAlienvault(domain string) []string {
	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", domain)
	fmt.Printf("[#] Fetching Subdomains from otx.alienvault.com for %s\n", domain)

	// Use gorequest to make HTTP request
	request := gorequest.New()
	_, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Printf("[!] Error fetching data from %s: %v\n", url, errs)
		return []string{}
	}

	// Parse JSON response
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)

	if err != nil {
		fmt.Printf("[!] Error decoding JSON for domain %s from %s\n", domain, url)
		return []string{}
	}

	subdomains := []string{}

	if passiveDNS, ok := data["passive_dns"].([]interface{}); ok {
		for _, entry := range passiveDNS {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				if hostname, ok := entryMap["hostname"].(string); ok {
					subdomains = append(subdomains, hostname)
				}
			}
		}
	}

	if len(subdomains) == 0 {
		fmt.Println("[X] No passive DNS data found.")
	}

	return subdomains
}

func fetchSubdomainsFromUrlscan(domain string) []string {
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=%s", domain)
	fmt.Printf("[#] Fetching Subdomains from urlscan.io for %s\n", domain)

	// Use gorequest to make HTTP request
	request := gorequest.New()
	_, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Printf("[!] Error fetching data from %s: %v\n", url, errs)
		return nil // Return an empty slice
	}

	// Parse JSON response
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)

	if err != nil {
		fmt.Printf("[!] Error decoding JSON for domain %s from %s\n", domain, url)
		return nil // Return an empty slice
	}

	subdomains := []string{}

	if results, ok := data["results"].([]interface{}); ok {
		for _, entry := range results {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				if domain, ok := entryMap["domain"].(string); ok {
					subdomains = append(subdomains, domain)
				}
			}
		}
	}

	return subdomains
}

func fetchSubdomainsFromAnubis(domain string) []string {
	url := fmt.Sprintf("https://jldc.me/anubis/subdomains/%s", domain)
	fmt.Printf("[#] Fetching Subdomains from Anubis for %s\n", domain)

	// Use gorequest to make HTTP request
	request := gorequest.New()
	_, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Printf("[!] Error fetching data from %s: %v\n", url, errs)
		return []string{}
	}

	// Parse JSON response
	var subdomains []string
	err := json.Unmarshal([]byte(body), &subdomains)

	if err != nil {
		fmt.Printf("Anubis response for %s is not in the expected format.\n", domain)
		return []string{}
	}

	return subdomains
}
func fetchSubdomainsFromHackerTarget(domain string) []string {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)
	fmt.Printf("[#] Fetching Subdomains from hackertarget.com for %s\n", domain)

	// Use gorequest to make HTTP request
	request := gorequest.New()
	_, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Printf("[!] Error fetching data from %s: %v\n", url, errs)
		return []string{}
	}

	// Parse CSV response
	scanner := bufio.NewScanner(strings.NewReader(body))
	var subdomains []string

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) > 0 {
			subdomains = append(subdomains, fields[0])
		}
	}

	if len(subdomains) == 0 {
		fmt.Println("[X] No subdomains found.")
	}

	return subdomains
}

func fetchSubdomainsFromRapiddns(domain string) []string {
	url := fmt.Sprintf("https://rapiddns.io/subdomain/%s?full=1#result", domain)
	fmt.Printf("[#] Fetching Subdomains from rapiddns.io for %s\n", domain)

	// Use gorequest to make HTTP request
	request := gorequest.New()
	_, body, errs := request.Get(url).End()

	if errs != nil {
		fmt.Printf("[!] Error fetching data from %s: %v\n", url, errs)
		return []string{}
	}

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		fmt.Printf("[X] Error parsing HTML for domain %s from %s: %v\n", domain, url, err)
		return []string{}
	}

	subdomains := []string{}
	websiteTable := doc.Find("table.table-striped").First()

	if websiteTable.Length() > 0 {
		websiteTable.Find("tbody").Each(func(_ int, tbody *goquery.Selection) {
			tbody.Find("tr").Each(func(_ int, tr *goquery.Selection) {
				subdomain := tr.Find("td").First().Text()
				subdomains = append(subdomains, subdomain)
			})
		})
	} else {
		fmt.Println("[X] No table element found on the page.")
	}

	return subdomains
}

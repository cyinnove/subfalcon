package sub88r

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

// Results holds the subdomain and wildcards results.
type Results struct {
	Subdomains []string
	Wildcards  []string
}

// Subber provides methods to scrape subdomains from various sources.
type Subber struct {
	Domain  string // Domain is the target domain for which subdomains will be scraped.
	Results *Results
}

// Getter method for retrieving subdomains
func (s *Subber) GetAllSubdomains() []string {
	return s.Results.Subdomains
}

// Getter method for retrieving wildcards
func (s *Subber) GetAllWildcards() []string {
	return s.Results.Wildcards
}

// RapidDNS scrapes subdomains from rapiddns.io, It returns a slice of subdomains and error.
func (s *Subber) RapidDNS() error {
	c := colly.NewCollector()
	c.OnHTML("tbody tr", func(h *colly.HTMLElement) {
		tdText := h.DOM.Find("td").First().Text()
		s.Results.Subdomains = append(s.Results.Subdomains, tdText)
	})

	url := fmt.Sprintf("https://rapiddns.io/subdomain/%s?full=1#result", s.Domain)
	c.Visit(url)

	return nil
}

// HackerTarget scrapes subdomains from hackertarget.com, It returns a slice of subdomains and error.
func (s *Subber) HackerTarget() error {
	// Send request
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", s.Domain)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Scrap subdomains
	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) > 1 {
			s.Results.Subdomains = append(s.Results.Subdomains, parts[0])
		}
	}

	return nil
}

// Anubis scrapes subdomains from Anubis via jldc.me, It returns a slice of subdomains and error.
func (s *Subber) Anubis() error {
	// Send Request
	url := fmt.Sprintf("https://jldc.me/anubis/subdomains/%s", s.Domain)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode json and scrap subdomains
	if err := json.NewDecoder(resp.Body).Decode(&s.Results.Subdomains); err != nil {
		return err
	}

	return nil
}

// urlscanResponse represents the response structure from urlscan.io API.
type urlscanResponse struct {
	Results []struct {
		Task struct {
			Domain string
		} `json:"task"`
	} `json:"results"`
}

// UrlScan scrapes subdomains from urlscan.io, It returns a slice of subdomains and error.
func (s *Subber) UrlScan() error {
	// Send Request
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=%s", s.Domain)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var result urlscanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Extract subdomains
	for _, res := range result.Results {
		s.Results.Subdomains = append(s.Results.Subdomains, res.Task.Domain)
	}

	return nil
}

// otxResults represents the response structure from Alien Vault (OTX) API.
type otxResults struct {
	PassiveDNS []struct {
		Hostname string `json:"hostname"`
	} `json:"passive_dns"`
}

// Otx scrapes subdomains from Alien Vault (OTX) via otx.alienvault.com, It returns a slice of subdomains and error.
func (s *Subber) Otx() error {
	// Send Request
	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/passive_dns", s.Domain)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode Json
	var res otxResults
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	// Scrap Subdomains
	for _, entry := range res.PassiveDNS {
		s.Results.Subdomains = append(s.Results.Subdomains, entry.Hostname)
	}

	return nil
}

// crtshResponse represents the response structure from crt.sh API.
type crtshResponse struct {
	NameValue string `json:"name_value"`
}

// CrtSh scrapes subdomains from crt.sh, It returns a slice of subdomains and slice of wildcards and error.
func (s *Subber) CrtSh() error {

	// Declare Response Structure
	var Responses []crtshResponse

	// Send Request
	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", s.Domain)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse JSON response
	if err := json.NewDecoder(resp.Body).Decode(&Responses); err != nil {
		return err
	}

	// Scrap subdomains and wildcards save them in a slices
	for _, response := range Responses {
		nameValue := response.NameValue
		if strings.Contains(nameValue, "\n") {
			subnameValues := strings.Split(nameValue, "\n")
			for _, subname := range subnameValues {
				subname = strings.TrimSpace(subname)
				if subname != "" {
					if strings.Contains(subname, "*") {
						s.Results.Wildcards = append(s.Results.Wildcards, subname)
					} else {
						s.Results.Subdomains = append(s.Results.Subdomains, subname)
					}
				}
			}
		} else {
			if strings.Contains(nameValue, "*") {
				s.Results.Wildcards = append(s.Results.Wildcards, nameValue)
			} else {
				s.Results.Subdomains = append(s.Results.Subdomains, nameValue)
			}
		}
	}

	return nil
}

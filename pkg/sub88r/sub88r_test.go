package sub88r

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllSubdomains(t *testing.T) {
	// Initialize Subber
	subber := Subber{
		Domain: "example.com",
		Results: &Results{
			Subdomains: []string{"sub1.example.com", "sub2.example.com"},
			Wildcards:  []string{"*.wildcard.example.com"},
		},
	}

	// Retrieve subdomains
	subdomains := subber.GetAllSubdomains()

	// Check if the length of retrieved subdomains matches the expected length
	expectedLength := 2
	if len(subdomains) != expectedLength {
		t.Errorf("Expected %d subdomains, got %d", expectedLength, len(subdomains))
	}
}

func TestGetAllWildcards(t *testing.T) {
	// Initialize Subber
	subber := Subber{
		Domain: "example.com",
		Results: &Results{
			Subdomains: []string{"sub1.example.com", "sub2.example.com"},
			Wildcards:  []string{"*.wildcard.example.com"},
		},
	}

	// Retrieve wildcards
	wildcards := subber.GetAllWildcards()

	// Check if the length of retrieved wildcards matches the expected length
	expectedLength := 1
	if len(wildcards) != expectedLength {
		t.Errorf("Expected %d wildcards, got %d", expectedLength, len(wildcards))
	}
}

func TestRapidDNS(t *testing.T) {
	t.Run("Valid Case", func(t *testing.T) {
		// Initialize Subber
		subber := Subber{
			Domain:  "example.com",
			Results: &Results{},
		}

		// Call RapidDNS
		err := subber.RapidDNS()
		assert.NoError(t, err, "RapidDNS failed")

		// Check if the results is correct
		expected := []string{"example.com", "www.example.com", "example.com", "www.example.com", "example.com", "example.com", "example.com"}
		actual := subber.Results.Subdomains

		assert.Equal(t, expected, actual)

	})
}

func TestHackerTarget(t *testing.T) {
	t.Run("Valid Case", func(t *testing.T) {
		// Initialize Subber
		subber := Subber{
			Domain:  "example.com",
			Results: &Results{},
		}

		// Call HackerTarget
		err := subber.HackerTarget()
		assert.NoError(t, err, "HackerTarget failed")

		// Define the expected subdomains
		expected := []string{"WWW.example.com", "Www.example.com", "www.example.com"}
		actual := subber.Results.Subdomains

		// Check if the results match the expected subdomains
		assert.Equal(t, expected, actual, "Unexpected subdomains")

	})
}

func TestAnubis(t *testing.T) {
	t.Run("Valid Case", func(t *testing.T) {
		// Initialize Subber
		subber := Subber{
			Domain:  "vulnweb.com",
			Results: &Results{},
		}

		// Call Anubis
		err := subber.Anubis()
		assert.NoError(t, err, "Anubis failed")

		// Define the expected subdomains
		expected := []string{"testhtml5.vulnweb.com", "testphp.vulnweb.com", "testasp.vulnweb.com", "testaspnet.vulnweb.com", "estphp.vulnweb.com", "antivirus1.vulnweb.com", "viruswall.vulnweb.com", "odincovo.vulnweb.com", "test.php.vulnweb.com", "tetphp.vulnweb.com", "virus.vulnweb.com", "rest.vulnweb.com"}
		actual := subber.Results.Subdomains

		// Check if the results match the expected subdomains
		assert.Equal(t, expected, actual, "Unexpected subdomains")

	})
}

func TestUrlScan(t *testing.T) {
	t.Run("Valid Case", func(t *testing.T) {
		// Initialize Subber
		subber := Subber{
			Domain:  "examplee.com",
			Results: &Results{},
		}

		// Call UrlScan
		err := subber.UrlScan()
		assert.NoError(t, err, "UrlScan failed")

		// Define the expected subdomains
		expected := []string{"vpn1.examplee.com"}
		actual := subber.Results.Subdomains

		// Check if the results match the expected subdomains
		assert.Equal(t, expected, actual, "Unexpected subdomains")

	})
}

func TestOtx(t *testing.T) {
	t.Run("successful case ", func(t *testing.T) {
		// Initialize Subber
		subber := Subber{
			Domain:  "testphp.vulnweb.com",
			Results: &Results{},
		}

		// Call Otx
		err := subber.Otx()
		assert.NoError(t, err, "Otx failed")

		// Define the expected subdomains
		expected := []string{
			"www.testphp.vulnweb.com",
			"testphp.vulnweb.com",
			"testphp.vulnweb.com",
			"testphp.vulnweb.com",
			"www.testphp.vulnweb.com",
			"testphp.vulnweb.com",
		}
		actual := subber.Results.Subdomains

		// Check if the results match the expected subdomains
		assert.Equal(t, expected, actual, "Unexpected subdomains")

	})
}

func TestCrtSh(t *testing.T) {
	t.Run("Fail Case", func(t *testing.T) {
		// Initialize Subber
		subber := Subber{
			Domain:  "testphp.vulnweb.com",
			Results: &Results{},
		}

		// Call CrtSh
		err := subber.CrtSh()
		println(subber.Results.Subdomains)
		assert.NoError(t, err, "CrtSh failed")

		// Define the expected subdomains
		expected := []string{}
		actual := subber.Results.Subdomains

		// Check if the results match the expected subdomains
		assert.Equal(t, expected, actual, "Unexpected subdomains")

	})
}

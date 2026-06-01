package fragscrapper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Website string

const (
	WebsiteFrag Website = "https://www.fragrantica.com"
	WebsiteParf Website = "https://www.parfumo.com"
)

type Scraper struct {
	website Website
	client  tls_client.HttpClient
	headers http.Header
}

type FragranceParams struct {
	FragranticaID int32
	Name          string
	Brand         string
	Country       string
	Gender        string
	RatingValue   string
	RatingCount   int32
	Year          int32
	TopNotes      string
	MiddleNotes   string
	BaseNotes     string
	Perfumer1     string
	Perfumer2     string
	Accord1       string
	Accord2       string
	Accord3       string
	Accord4       string
	Accord5       string
	Accord6       string
	Accord7       string
	Accord8       string
	Accord9       string
	Accord10      string
}

// HTTP client that looks like a real browser to avoid being blocked
func newScraper(website Website) (*Scraper, error) {
	jar := tls_client.NewCookieJar()
	client, err := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithCookieJar(jar),
		tls_client.WithRandomTLSExtensionOrder(),
	)
	if err != nil {
		return nil, err
	}

	headers := http.Header{
		"accept": {
			"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		},
		"accept-encoding": {"gzip, deflate, br"},
		"accept-language": {"en-US,en;q=0.9"},
		"cache-control":   {"max-age=0"},
		"sec-ch-ua": {
			`"Chromium";v="120", "Not(A:Brand";v="24", "Google Chrome";v="120"`,
		},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {`"Windows"`},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent": {
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		http.HeaderOrderKey: {
			"accept",
			"accept-encoding",
			"accept-language",
			"cache-control",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"sec-fetch-user",
			"upgrade-insecure-requests",
			"user-agent",
		},
	}

	return &Scraper{website: website, client: client, headers: headers}, nil
}

// Get html from url
func (s *Scraper) getPageBody(url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = s.headers

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

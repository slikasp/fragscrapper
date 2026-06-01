package fragscrapper

type ParfScraper struct {
	*Scraper
}

func NewParfScraper() (*ParfScraper, error) {
	base, err := newScraper(WebsiteParf)
	if err != nil {
		return nil, err
	}

	return &ParfScraper{
		Scraper: base,
	}, nil
}

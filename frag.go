package fragscrapper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FragScraper struct {
	*Scraper
}

func NewFragScraper() (*FragScraper, error) {
	base, err := newScraper(WebsiteFrag)
	if err != nil {
		return nil, err
	}

	return &FragScraper{
		Scraper: base,
	}, nil
}

// Get all details about a fragrance from the frag******a website
func (s *FragScraper) GetFragranceParams(url string) (FragranceParams, error) {
	params := FragranceParams{}

	doc, err := s.getPageBody(url)
	if err != nil {
		return params, fmt.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		return params, fmt.Errorf("Empty response body")
	}

	// Gender - for men / for women / for women and men (or vice versa?)
	gender := getGender(doc)
	params.Gender = gender

	// RatingValue & RatingCount
	rVal, rCount := getRatings(doc)
	if rVal != "" {
		params.RatingValue = rVal
		params.RatingCount = int32(rCount)
	}

	// Year
	year, known := getYear(doc)
	if known {
		params.Year = int32(year)
	}

	// Notes [topNotes, middleNotes, baseNotes]
	notes := getNotes(doc)
	params.TopNotes = notes[0]
	params.MiddleNotes = notes[1]
	params.BaseNotes = notes[2]

	// Parfumer1-2 (["unknown"] is received if there were none found)
	perfumers := getPerfumers(doc)
	params.Perfumer1 = perfumers[0]
	if len(perfumers) > 1 {
		params.Perfumer2 = perfumers[1]
	}

	// unpack Accords1-10
	accords := getAccords(doc)
	for i, v := range accords {
		if i >= 10 {
			break
		}

		switch i {
		case 0:
			params.Accord1 = v
		case 1:
			params.Accord2 = v
		case 2:
			params.Accord3 = v
		case 3:
			params.Accord4 = v
		case 4:
			params.Accord5 = v
		case 5:
			params.Accord6 = v
		case 6:
			params.Accord7 = v
		case 7:
			params.Accord8 = v
		case 8:
			params.Accord9 = v
		case 9:
			params.Accord10 = v
		}
	}

	return params, nil
}

// Get perfumes from user's wardrobe
func (s *FragScraper) GetMemberWardrobe(id int) (map[string][]string, error) {
	fragrances := make(map[string][]string)

	shelves := make([]*goquery.Selection, 10)

	wardrobeUrl := fmt.Sprintf("%s/member/%d#wardrobe", s.website, id)
	customUrl := fmt.Sprintf("%s/member/%d#custom", s.website, id)

	wardrobeDoc, err := s.getPageBody(wardrobeUrl)
	if err != nil {
		return fragrances, fmt.Errorf("Read wardrobe body failed: %s", err)
	}
	if wardrobeDoc == nil {
		return fragrances, fmt.Errorf("Empty response body")
	}
	customDoc, err := s.getPageBody(customUrl)
	if err != nil {
		return fragrances, fmt.Errorf("Read custom body failed: %s", err)
	}
	if customDoc == nil {
		return fragrances, fmt.Errorf("Empty response body")
	}

	wardrobeDoc.Find(`div[role="tabpanel"][data-headlessui-state="selected"]`).Find("div").Each(func(i int, item *goquery.Selection) {
		shelves = append(shelves, item)
	})

	customDoc.Find(`div[role="tabpanel"][data-headlessui-state="selected"]`).Find("div").Each(func(i int, item *goquery.Selection) {
		shelves = append(shelves, item)
	})

	return fragrances, nil
}

// Get the perfumer country from the frag******a page
func (s *FragScraper) getPerfumerCountry(perfumer string) (string, error) {
	url := fmt.Sprintf("%s/%s.html", s.website, perfumer)

	doc, err := s.getPageBody(url)
	if err != nil {
		return "", fmt.Errorf("Read body failed: %s", err)
	}

	if doc == nil {
		return "", fmt.Errorf("Empty response body")
	}

	var results []string

	doc.Find("main div.col-span-8.col-start-5.md\\:col-span-full").Find("a").Each(func(i int, item *goquery.Selection) {
		results = append(results, strings.TrimSpace(item.Text()))
	})

	if len(results) == 0 {
		return "", fmt.Errorf("No perfumer details in %s", url)
	}

	return results[0], nil
}

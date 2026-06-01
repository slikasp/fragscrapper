package fragscrapper

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getGender(doc *goquery.Document) string {
	sex := strings.ToLower(strings.TrimSpace(doc.Find("#toptop").Find("span").Text()))

	var gender string

	switch sex {
	case "for men":
		gender = "men"
	case "for women":
		gender = "women"
	default:
		gender = "unisex"
	}

	return gender
}

func getRatings(doc *goquery.Document) (string, int) {
	rating := doc.Find(`[itemprop="aggregateRating"]`)

	value := rating.Find(`[itemprop="ratingValue"]`).Text()

	countSel := rating.Find(`[itemprop="ratingCount"]`)
	countStr, exists := countSel.Attr("content")
	if !exists {
		return "", 0
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", 0
	}

	return value, count
}

func getYear(doc *goquery.Document) (int, bool) {
	s := doc.Find("head > title").Text()

	parts := strings.Fields(s)
	if len(parts) == 0 {
		return 0, false
	}

	last := parts[len(parts)-1]

	year, err := strconv.Atoi(last)
	if err != nil {
		return 0, false
	}

	//basic sanity check
	if year < 1000 || year > 9999 {
		return 0, false
	}

	return year, true
}

func getNotes(doc *goquery.Document) []string {
	var results []string

	// TopNotes
	// #pyramid > div.mx-auto.max-w-md
	var topNotes []string
	doc.Find("#pyramid div.mx-auto.max-w-md").Find("div.flex.flex-wrap.justify-center.items-end.py-3.px-2.pyramid-level-container").Find("span").Each(func(i int, item *goquery.Selection) {
		topNotes = append(topNotes, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	// MiddleNotes
	// #pyramid > div.mx-auto.max-w-xl
	var middleNotes []string
	doc.Find("#pyramid div.mx-auto.max-w-xl").Find("div.flex.flex-wrap.justify-center.items-end.py-3.px-2.pyramid-level-container").Find("span").Each(func(i int, item *goquery.Selection) {
		middleNotes = append(middleNotes, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	// BaseNotes
	// #pyramid > div.mx-auto.max-w-2xl
	var baseNotes []string
	doc.Find("#pyramid div.mx-auto.max-w-2xl").Find("div.flex.flex-wrap.justify-center.items-end.py-3.px-2.pyramid-level-container").Find("span").Each(func(i int, item *goquery.Selection) {
		baseNotes = append(baseNotes, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	results = append(results, strings.Join(topNotes, ", "))
	results = append(results, strings.Join(middleNotes, ", "))
	results = append(results, strings.Join(baseNotes, ", "))

	return results
}

func getPerfumers(doc *goquery.Document) []string {
	var results []string

	// Perfumer1 (optional)
	// Perfumer2 (optional)
	// website has space for 1-4 perfumers, only need 2, set Perfumer1 to unknown if none found
	doc.Find("main div.grid.grid-cols-2.md\\:grid-cols-3.lg\\:grid-cols-4.gap-3.md\\:gap-4").Find("span").Each(func(i int, item *goquery.Selection) {
		results = append(results, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	if len(results) == 0 {
		results = append(results, "unknown")
	}

	return results
}

func getAccords(doc *goquery.Document) []string {
	var results []string

	doc.Find("main div.flex.flex-col.w-full.max-w-\\[280px\\].md\\:max-w-\\[320px\\]").Find("span.truncate").Each(func(i int, item *goquery.Selection) {
		results = append(results, strings.ToLower(strings.TrimSpace(item.Text())))
	})

	return results
}

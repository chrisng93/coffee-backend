package data

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/chrisng93/coffee-backend/models"

	"github.com/chrisng93/coffee-backend/db"
	"github.com/gocolly/colly"
)

// yelpAttribute is sent into the attribute channel when we scrape each Yelp URL. It contains
// information about whether a Yelp business is good for working and/or has free wifi.
type yelpAttributes struct {
	yelpURL          string
	isGoodForWorking bool
	hasWifi          bool
}

func scrapeCoffeeShopYelpURLs(databaseOps *db.DatabaseOps) error {
	// TODO: Get all coffee shops, scrape their Yelp URL to see if study friendly.
	// Get all coffee shops (db call)
	// For each coffee shop, scrape their Yelp URL. Get info on "good for studying" and "wifi"
	// Update coffee shops (db call)

	coffeeShops, err := databaseOps.GetCoffeeShops()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(coffeeShops))
	// Use attributeChan to facilitate communication with goroutines used for scraping.
	attributeChan := make(chan *yelpAttributes, len(coffeeShops))
	// Create mapping of Yelp URL to pointer to coffee shop. Use this mapping later for getting
	// the coffee shop associated with a Yelp URL.
	yelpURLToCoffeeShop := map[string]*models.CoffeeShop{}
	for _, coffeeShop := range coffeeShops {
		go scrapeYelpURL(wg, coffeeShop.YelpURL, attributeChan)
		// Pause for 500ms between requests so we don't bombard Yelp.
		time.Sleep(500 * time.Millisecond)
		yelpURLToCoffeeShop[coffeeShop.YelpURL] = coffeeShop
	}

	wg.Wait()
	close(attributeChan)

	// We only want to update coffee shops that have an IsGoodForStudying attribute that has
	// changed.
	var changedCoffeeShops []*models.CoffeeShop
	for yelpAttribute := range attributeChan {
		coffeeShop := yelpURLToCoffeeShop[yelpAttribute.yelpURL]
		newIsGoodForStudying := yelpAttribute.isGoodForWorking && yelpAttribute.hasWifi
		if coffeeShop.IsGoodForStudying != newIsGoodForStudying {
			coffeeShop.IsGoodForStudying = yelpAttribute.isGoodForWorking && yelpAttribute.hasWifi
			changedCoffeeShops = append(changedCoffeeShops, coffeeShop)
		}
	}

	return databaseOps.UpdateCoffeeShops(changedCoffeeShops)
}

func scrapeYelpURL(wg *sync.WaitGroup, yelpURL string, attributeChan chan *yelpAttributes) {
	c := colly.NewCollector()

	c.OnHTML(".ywidget .ylist .short-def-list", func(listElem *colly.HTMLElement) {
		// Need to defer wg.Done() here and in onError because one of these two must be called by
		// colly.Collector.
		defer wg.Done()
		var isGoodForWorking, hasWifi bool
		listElem.ForEach("dl", func(_ int, listItemElem *colly.HTMLElement) {
			// Check Yelp's "More business info" for the "Good for Working" and "Wi-fi" attributes.
			formattedText := formatString(listItemElem.Text)
			if strings.Contains(formattedText, "goodforworking") {
				isGoodForWorking = attributeHasPositiveValue(formattedText, "goodforworking", "yes")
			}
			if strings.Contains(formattedText, "wi-fi") {
				hasWifi = attributeHasPositiveValue(formattedText, "wi-fi", "free")
			}
		})
		attributeChan <- &yelpAttributes{
			yelpURL:          yelpURL,
			isGoodForWorking: isGoodForWorking,
			hasWifi:          hasWifi,
		}
	})

	c.Visit(yelpURL)

	c.OnError(func(_ *colly.Response, err error) {
		defer wg.Done()
		log.Printf("Error scraping Yelp URL: %v", err)
	})
}

func attributeHasPositiveValue(attributeAndValue string, attribute string, positiveCase string) bool {
	return strings.SplitAfter(attributeAndValue, attribute)[1] == positiveCase
}

// formatString takes a list item element from Yelp's "More business info" section and puts it in
// the following format: ${attribute}${value}. Ex: goodforworkingyes.
func formatString(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Replace(text, " ", "", -1)
	text = strings.Replace(text, "\n", "", -1)
	return strings.ToLower(text)
}

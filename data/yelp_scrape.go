package data

import (
	"fmt"
	"strings"

	"github.com/chrisng93/coffee-backend/db"
	"github.com/gocolly/colly"
)

func scrapeCoffeeShopYelpURLs(databaseOps *db.DatabaseOps) error {
	// TODO: Get all coffee shops, scrape their Yelp URL to see if study friendly.
	// Get all coffee shops (db call)
	// For each coffee shop, scrape their Yelp URL. Get info on "good for studying" and "wifi"
	// Update coffee shops (db call)

	coffeeShops, err := databaseOps.GetCoffeeShops()
	if err != nil {
		return err
	}

	for _, coffeeShop := range coffeeShops {
		fmt.Println(coffeeShop)
		go scrapeYelpURL(coffeeShop.YelpURL)
	}

	return nil
}

func scrapeYelpURL(yelpURL string) {
	// Instantiate default collector
	c := colly.NewCollector()

	var found bool
	c.OnHTML(".ywidget .ylist .short-def-list", func(listElem *colly.HTMLElement) {
		found = true
		listElem.ForEach("dl", func(i int, listItemElem *colly.HTMLElement) {
			parseListItem(listItemElem, "goodforworking")
			parseListItem(listItemElem, "wi-fi")
		})
	})
	if !found {
		// fmt.Println(yelpURL)
	}

	c.Visit(yelpURL)
}

func parseListItem(listItemElem *colly.HTMLElement, attribute string) {
	str := strings.ToLower(strings.Replace(strings.Replace(strings.TrimSpace(listItemElem.Text), " ", "", -1), "\n", "", -1))
	if len(str) > len(attribute) && str[:len(attribute)] == attribute {
		// fmt.Println(strings.SplitAfter(str, attribute))
	}
}

package data

import (
	"log"
	"sync"
	"time"

	"github.com/chrisng93/coffee-backend/clients/yelp"
	"github.com/chrisng93/coffee-backend/db"
	"github.com/chrisng93/coffee-backend/models"
)

// ExcludedCoffeeShops is a list of coffee shops not to include.
var ExcludedCoffeeShops = []string{"Starbucks", "Dunkin' Donuts"}

// Neighborhoods defines which neighboords to look at coffee shops in.
var Neighborhoods = []string{"Lower Manhattan", "Williamsburg", "Greenpoint"}

func getAndUpsertYelpData(databaseOps *db.DatabaseOps, yelpClient *yelp.Client) error {
	// TODO: We have 5000 calls to the Yelp API every day. Right now, these two calls hit their
	// API ~100x, and each call gets more than the max amount (1000) of results from Yelp. Think
	// about ways to better target certain areas of NYC/Brooklyn where there might be higher
	// density of coffee shops that don't show up in these queries. Potentially search different
	// neighborhoods rather than boroughs.
	log.Println("Calling Yelp API to get coffee shop data...")
	var queries []*yelp.SearchBusinessesParams
	for _, neighborhood := range Neighborhoods {
		queries = append(queries, &yelp.SearchBusinessesParams{
			Location:   neighborhood,
			SearchTerm: "best coffee shops",
			Categories: "coffee,coffeeroasteries",
		})
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(queries))
	coffeeShopChan := make(chan []*yelp.Business, len(queries))
	for _, searchParams := range queries {
		time.Sleep(10 * time.Second)
		go getCoffeeShops(wg, coffeeShopChan, yelpClient, searchParams)
	}

	wg.Wait()
	close(coffeeShopChan)

	var filteredCoffeeShopsYelp []*yelp.Business
	for coffeeShops := range coffeeShopChan {
		filteredCoffeeShopsYelp = append(filteredCoffeeShopsYelp, filterCoffeeShops(coffeeShops)...)
	}

	var coffeeShops []*models.CoffeeShop
	// We use this to detect duplicates between the queries.
	yelpIDToSeen := map[string]bool{}
	for _, yelpCoffeeShop := range filteredCoffeeShopsYelp {
		// If we've already seen this coffee shop (duplicate), then don't add it again.
		if _, ok := yelpIDToSeen[yelpCoffeeShop.YelpID]; ok {
			continue
		}

		coffeeShops = append(coffeeShops, &models.CoffeeShop{
			Name: yelpCoffeeShop.Name,
			Coordinates: &models.Coordinates{
				Latitude:  yelpCoffeeShop.Coordinates.Latitude,
				Longitude: yelpCoffeeShop.Coordinates.Longitude,
			},
			YelpID:  yelpCoffeeShop.YelpID,
			YelpURL: yelpCoffeeShop.URL,
		})
		yelpIDToSeen[yelpCoffeeShop.YelpID] = true
	}

	log.Printf("Processed %v coffee shops from Yelp...", len(filteredCoffeeShopsYelp))
	return databaseOps.InsertOrUpdateCoffeeShops(coffeeShops)
}

func getCoffeeShops(
	wg *sync.WaitGroup,
	coffeeShopChan chan []*yelp.Business,
	yelpClient *yelp.Client,
	params *yelp.SearchBusinessesParams,
) {
	defer wg.Done()
	bestCoffeeShopsYelp, err := yelpClient.SearchBusinesses(params)
	if err != nil {
		log.Printf("Error getting data from Yelp: %v", err)
	}
	coffeeShopChan <- bestCoffeeShopsYelp
}

func filterCoffeeShops(coffeeShops []*yelp.Business) []*yelp.Business {
	var filteredCoffeeShops []*yelp.Business
	for _, coffeeShop := range coffeeShops {
		if includeCoffeeShop(coffeeShop) {
			filteredCoffeeShops = append(filteredCoffeeShops, coffeeShop)
		}
	}
	return filteredCoffeeShops
}

func includeCoffeeShop(coffeeShop *yelp.Business) bool {
	// Only NY for now.
	if coffeeShop.Location.State != "NY" {
		return false
	}
	for _, excludedCoffeeShop := range ExcludedCoffeeShops {
		// Don't inlcude coffee shop if it appears in the excluded coffee shop list.
		if coffeeShop.Name == excludedCoffeeShop {
			return false
		}
	}
	// Most coffee shops that we're interested in don't have more than 3 categories - they're more
	// refined in their foucs.
	if len(coffeeShop.Categories) > 3 {
		return false
	}
	for _, category := range coffeeShop.Categories {
		// Random ice cream shops show up as coffee shops as well.
		if category.Alias == "icecream" {
			return false
		}
	}
	// TODO: Do a better job of this. This is very unrefined.
	return (coffeeShop.Rating > 4 && coffeeShop.ReviewCount > 50) || coffeeShop.ReviewCount > 200
}

package data

import (
	"log"

	"github.com/chrisng93/coffee-backend/db"
	"github.com/chrisng93/coffee-backend/yelp"
)

// InitializeCronJobs initializes cron jobs to get data from Yelp and Instagram.
func InitializeCronJobs(databaseOps *db.DatabaseOps, yelpClient *yelp.Client) {
	// TODO: Create cron job to call Yelp's API every day to businesses.
	err := getAndUpsertYelpData(databaseOps, yelpClient)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("Successfully upserted Yelp data.")
	}

	// TODO: Create cron job to scrape Yelp's businesses.
	// err := scrapeCoffeeShopYelpURLs(databaseOps)
	// if err != nil {
	// 	log.Print(err)
	// }

	// TODO: Create cron job to call Instagram's API.
}

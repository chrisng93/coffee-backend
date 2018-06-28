package data

import (
	"log"

	"github.com/robfig/cron"

	"github.com/chrisng93/coffee-backend/clients/yelp"
	"github.com/chrisng93/coffee-backend/db"
)

// InitializeCronJobs initializes cron jobs to get data from Yelp and Instagram.
func InitializeCronJobs(databaseOps *db.DatabaseOps, yelpClient *yelp.Client) {
	c := cron.New()
	c.AddFunc("@daily", func() {
		err := getAndUpsertYelpData(databaseOps, yelpClient)
		if err != nil {
			log.Print(err)
		} else {
			log.Print("Successfully upserted Yelp data.")
		}
	})
	c.AddFunc("@daily", func() {
		err := scrapeCoffeeShopYelpURLs(databaseOps)
		if err != nil {
			log.Print(err)
		} else {
			log.Printf("Successfully scraped Yelp URLs")
		}
	})
	c.Start()

	// TODO: Create cron job to call Instagram's API.

	// Populate database on startup.
	// err := getAndUpsertYelpData(databaseOps, yelpClient)
	// if err != nil {
	// 	log.Print(err)
	// } else {
	// 	log.Print("Successfully upserted Yelp data.")
	// }
	// err = scrapeCoffeeShopYelpURLs(databaseOps)
	// if err != nil {
	// 	log.Print(err)
	// } else {
	// 	log.Printf("Successfully scraped Yelp URLs")
	// }
}

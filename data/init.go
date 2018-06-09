package data

import (
	"database/sql"
	"log"

	"github.com/chrisng93/coffee-backend/yelp"
)

// InitializeCronJobs initializes cron jobs to get data from Yelp and Instagram.
func InitializeCronJobs(db *sql.DB, yelpClient *yelp.Client) {
	// TODO: Create cron job to call Yelp's API every day to businesses.
	err := getAndUpsertYelpData(db, yelpClient)
	if err != nil {
		log.Print(err)
	} else {
		log.Print("Successfully upserted Yelp data.")
	}
	// TODO: Create cron job to call Instagram's API.
}

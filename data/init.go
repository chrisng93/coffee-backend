package data

import (
	"database/sql"

	"github.com/chrisng93/coffee-backend/yelp"
)

// InitializeCronJobs initializes cron jobs to get data from Yelp and Instagram.
func InitializeCronJobs(db *sql.DB, yelpClient *yelp.Client) {
	// TODO: Create cron job to call Yelp's API every day to businesses.
	getYelpData(yelpClient)
	// TODO: Create cron job to call Instagram's API.
}

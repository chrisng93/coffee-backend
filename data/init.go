package data

import (
	"github.com/chrisng93/coffee-backend/yelp"
)

// InitializeCronJobs initializes cron jobs to get data from Yelp and Instagram.
func InitializeCronJobs(yelpClient *yelp.Client) {
	// TODO: Get rid of this - it's just for testing.
	// TODO: Create cron job to call Yelp's API every day to businesses.
	getYelpData(yelpClient)
	// TODO: Create cron job to call Instagram's API.
}

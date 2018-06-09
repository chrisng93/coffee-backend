package db

import (
	"database/sql"
	"fmt"

	"github.com/chrisng93/coffee-backend/yelp"
)

// InsertOrUpdateYelpData inserts or updates filtered Yelp data from the /v3/businesses/search
// endpoint.
func InsertOrUpdateYelpData(db *sql.DB, coffeeShops []*yelp.Business) error {
	var failedTransactions []string

	err := createTransaction(db, func(tx *sql.Tx) error {
		for _, coffeeShop := range coffeeShops {
			// TODO: Right now, there is a separate statement for each upsert. This is
			// time-intensive. Change this to be a bulk upsert.
			_, err := tx.Exec(
				`INSERT INTO coffeeshop.shop (shop_name, lat, lng, yelp_id, yelp_url)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (yelp_id)
				 DO UPDATE SET (shop_name, lat, lng, yelp_url) = ($1, $2, $3, $5)`,
				coffeeShop.Name,
				coffeeShop.Coordinates.Latitude,
				coffeeShop.Coordinates.Longitude,
				coffeeShop.YelpID,
				coffeeShop.URL,
			)
			if err != nil {
				failedTransactions = append(failedTransactions, coffeeShop.YelpID)
				fmt.Println("unsuccessful db op for", coffeeShop.Name)
			} else {
				fmt.Println("successfully inserted", coffeeShop.Name)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(failedTransactions) > 0 {
		return fmt.Errorf("Error inserting or updating coffee shops: %v", failedTransactions)
	}
	return nil
}

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
	// TODO: Generalize transaction code.
	tx, err := db.Begin()
	var success bool
	defer func() {
		if success {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err != nil {
		return fmt.Errorf("Error creating transaction: %v", err)
	}

	for _, coffeeShop := range coffeeShops {
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

	if len(failedTransactions) > 0 {
		return fmt.Errorf("Error inserting or updating coffee shops: %v", failedTransactions)
	}
	return nil
}

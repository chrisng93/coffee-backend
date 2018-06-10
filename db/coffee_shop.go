package db

import (
	"database/sql"
	"fmt"

	"github.com/chrisng93/coffee-backend/models"
)

// InsertOrUpdateCoffeeShops inserts or updates coffee shops.
func (ops *DatabaseOps) InsertOrUpdateCoffeeShops(coffeeShops []*models.CoffeeShop) error {
	var failedTransactions []string

	err := createTransaction(ops.db, func(tx *sql.Tx) error {
		for _, coffeeShop := range coffeeShops {
			// TODO: Right now, there is a separate statement for each upsert. This is
			// time-intensive. Change this to be a bulk upsert.
			_, err := tx.Exec(
				`INSERT INTO coffeeshop.shop (name, lat, lng, yelp_id, yelp_url)
				 VALUES ($1, $2, $3, $4, $5)
				 ON CONFLICT (yelp_id)
				 DO UPDATE SET (name, lat, lng, yelp_url) = ($1, $2, $3, $5)`,
				coffeeShop.Name,
				coffeeShop.Coordinates.Latitude,
				coffeeShop.Coordinates.Longitude,
				coffeeShop.YelpID,
				coffeeShop.YelpURL,
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

// GetCoffeeShops gets all of the coffee shops in the database.
func (ops *DatabaseOps) GetCoffeeShops() ([]*models.CoffeeShop, error) {
	var coffeeShops []*models.CoffeeShop
	err := createTransaction(ops.db, func(tx *sql.Tx) error {
		rows, err := tx.Query(`
			SELECT id, last_updated, name, lat, lng, yelp_id, yelp_url, has_good_coffee,
				   is_good_for_studying
			FROM coffeeshop.shop
		`)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			coffeeShop := models.CoffeeShop{}
			coffeeShop.Coordinates = &models.Coordinates{}
			err := rows.Scan(
				&coffeeShop.ID,
				&coffeeShop.LastUpdated,
				&coffeeShop.Name,
				&coffeeShop.Coordinates.Latitude,
				&coffeeShop.Coordinates.Longitude,
				&coffeeShop.YelpID,
				&coffeeShop.YelpURL,
				&coffeeShop.HasGoodCoffee,
				&coffeeShop.IsGoodForStudying,
			)
			if err != nil {
				return err
			}
			coffeeShops = append(coffeeShops, &coffeeShop)
		}
		err = rows.Err()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return coffeeShops, nil
}

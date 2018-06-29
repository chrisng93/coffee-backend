# Cofee Around Me - Backend

## Setting up your technologies.
Make sure you have Go and Postgres installed.

## Setting up Postgres.
Run the following command from the coffee-backend root to create the coffeeshop schema:
```
psql -U ${db_user} -d ${db_name} -f ./schema/create_coffeeshop_schema.sql
```

## Running the API.
1. Create an app with [Yelp](https://www.yelp.com/fusion) to get your client ID and API key.
1. Register an app with [Google Maps](https://developers.google.com/maps/documentation/) to get your API key(s).
1. Run the following command from the root directory to install all of the third-party dependencies:
```
dep ensure
```
1. Run the following command from the root directory to start the server:
```
go run main.go --port=${port} --db_user=${db_user} --db_password=${db_password} --db_name=${db_name} --yelp_client_id=${yelp_client_id} --yelp_api_key=${yelp_api_key} --google_maps_api_key=${google_maps_api_key}
```

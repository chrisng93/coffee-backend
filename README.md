# coffee-backend

## Setting up Database.
Run the following command from the coffee-backend root to create the coffeeshop schema:
```
psql -U ${db_user} -d ${db_name} -f ./schema/create_coffeeshop_schema.sql
```

## Running the API.
1. Create an app with [Yelp](https://www.yelp.com/fusion) to get your client ID and API key.
1. TODO: Dependencies.
1. Run the following command from the coffee-backend root to start the server:
```
go run main.go --yelp_client_id=${yelp_client_id} --yelp_api_key=${yelp_api_key} --db_user=${db_user} --db_password=${db_password}
```

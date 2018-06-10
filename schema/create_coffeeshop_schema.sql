DROP SCHEMA IF EXISTS coffeeshop CASCADE;
CREATE SCHEMA coffeeshop;
SET SEARCH_PATH TO coffeeshop;

CREATE TABLE shop (
    -- Auto-generated unique ID.
    id BIGSERIAL PRIMARY KEY,
    -- Timestamp of last updated time.
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Coffee shop's name, as determined by Yelp.
    name TEXT NOT NULL,
    -- Shop's latitude.
    lat DECIMAL NOT NULL,
    -- Shop's longitude.
    lng DECIMAL NOT NULL,
    -- Shop's Yelp ID.
    yelp_id TEXT UNIQUE NOT NULL,
    -- Shop's Yelp URL.
    yelp_url TEXT UNIQUE NOT NULL,
    -- Whether or not the shop has good coffee.
    -- TODO: More on how this is determined.
    has_good_coffee BOOLEAN NOT NULL DEFAULT false,
    -- Whether or not the shop is good for studying. Determined by scraping Yelp URL and checking
    -- if the shop has wifi and is "good for working" (determined by Yelp).
    is_good_for_studying BOOLEAN NOT NULL DEFAULT false
    -- TODO: Uncomment these when add Instagram.
    -- is_instagrammable BOOLEAN NOT NULL DEFAULT false,
    -- instagram_id TEXT
);

CREATE OR REPLACE FUNCTION update_last_updated()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER update_shop_last_updated BEFORE UPDATE ON shop FOR EACH ROW EXECUTE PROCEDURE update_last_updated();

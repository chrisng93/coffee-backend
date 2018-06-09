DROP SCHEMA IF EXISTS coffeeshop CASCADE;
CREATE SCHEMA coffeeshop;
SET SEARCH_PATH TO coffeeshop;

CREATE TABLE shop (
    -- Auto-generated unique ID.
    id BIGINT PRIMARY KEY NOT NULL,
    -- Coffee shop's name, as determined by Yelp.
    shop_name TEXT NOT NULL,
    -- Shop's latitude.
    lat DECIMAL NOT NULL,
    -- Shop's longitude.
    lng DECIMAL NOT NULL,
    -- Shop's Yelp ID.
    yelp_id TEXT NOT NULL,
    -- Shop's Yelp URL.
    yelp_url TEXT NOT NULL,
    -- Whether or not the shop has good coffee.
    -- TODO: More on how this is determined.
    has_good_coffee BOOLEAN NOT NULL DEFAULT false,
    -- Whether or not the shop is good for studying. Determined by scraping Yelp URL and checking
    -- if the shop has wifi and is "good for working" (determined by Yelp).
    is_good_for_studying BOOLEAN NOT NULL DEFAULT false
    -- TODO: Uncomment these when add Instagram.
    -- is_instagrammable BOOLEAN DEFAULT false,
    -- instagram_id TEXT
);

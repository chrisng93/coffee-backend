# Use current version of go
FROM golang:latest 

MAINTAINER Chris Ng, chris.ng93@gmail.com

# Create app directory
RUN mkdir -p /go/src/github.com/chrisng93/coffee-backend
WORKDIR /go/src/github.com/chrisng93/coffee-backend

COPY . .
RUN dep ensure

# Build and run app
RUN go build -o main .
CMD ./main --port=8080 --db_user=postgres --db_password=postgres --db_host=postgres --yelp_client_id=yEWd7aXEfny_E9-8M2rG0Q --yelp_api_key=SIazOV9BrajOHV_N4d14xl9NYmglJJGKTf78mJeXwrdaPWXmvNejfMyWY_V8b7eJcktoRshznEWnptIUk2QSpuu_ptFP8IRXK-hgj3VJqQhv7hdbh2q0wV5e6dYbW3Yx --google_maps_api_key=AIzaSyANRSxSZBkrzse-7sv9WktRNqorthgKy2M,AIzaSyAawkgM6ykbI55b_lcWv6_cH1UxgonzA-s,AIzaSyALPkffaQAtq6Apo1CCN07xQQGu9kaqzB4,AIzaSyCQu8zSvsULxmsdpb5mc3ybtmA7MAGMfxM

# Use current version of go
FROM golang:latest 

# Install dependency tool
RUN go get github.com/golang/dep && go install github.com/golang/dep/cmd/dep

# Create app directory
RUN mkdir -p /go/src/github.com/chrisng93/coffee-backend
WORKDIR /go/src/github.com/chrisng93/coffee-backend

COPY . .
RUN dep ensure

# Build and run app
RUN go build -o main .
CMD ./main --port=80 --db_user=${db_user} --db_password=${db_password} --db_host=${db_host} --yelp_client_id=${yelp_client_id} --yelp_api_key=${yelp_api_key} --google_maps_api_key=${google_maps_api_key}

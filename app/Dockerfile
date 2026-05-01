# syntax=docker/dockerfile:1.23

# Build the application from source
FROM golang:1.26-alpine AS build-stage

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /stormsensor .

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static:nonroot AS build-release-stage

WORKDIR /

COPY --from=build-stage /stormsensor /stormsensor

USER nonroot:nonroot

ENTRYPOINT ["/stormsensor", "/var/lib/storm-to-mqtt-gw/config.json"]

# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /edgeswitch-to-mqtt-gw

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /edgeswitch-to-mqtt-gw /edgeswitch-to-mqtt-gw

USER nonroot:nonroot

ENTRYPOINT ["/edgeswitch-to-mqtt-gw", "/var/lib/edgeswitch-to-mqtt-gw/config.json"]

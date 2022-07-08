# syntax=docker/dockerfile:1

FROM golang:1.18-bullseye as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -v -o /traces-test-image

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /traces-test-image /traces-test-image
USER nonroot:nonroot
ENTRYPOINT ["/traces-test-image"]

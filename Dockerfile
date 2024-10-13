# Build image with all the build tools included
FROM golang:1.22.7 as builder

RUN mkdir /build

COPY ./ /build

RUN go build -C /build -o /opt/stock-ticker cmd/stock-ticker.go

# Smaller deployment image with no build tools included
FROM ubuntu

COPY --from=builder /opt/stock-ticker /opt/stock-ticker

# Base image doesn't have any CA certificates installed, so adding all the Google Trust Store CAs as that's what signs the stock ticker certificates.
COPY certificate-authorites/ca.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/opt/stock-ticker"]
FROM golang:1.14-alpine AS builder

WORKDIR /go/src/db-transaction-example

# We want to use `make install`
RUN apk add --update bash make

# Leverage docker cache to prevent dependencies download if not changed
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source code to be compiled
COPY . .

# Compile and install the application in /bin
RUN make install

# ---

FROM alpine

COPY --from=builder /bin/app /bin/app

EXPOSE 8080
CMD ["/bin/app"]

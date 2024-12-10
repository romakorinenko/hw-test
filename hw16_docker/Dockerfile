FROM golang:1.20-alpine AS build

WORKDIR /build

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./bin/app cmd/main.go

FROM alpine:3.21.0

COPY --from=build /build/bin/app /

CMD ["/app"]

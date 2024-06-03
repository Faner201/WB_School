FROM golang:1.21.4-alpine3.18 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . /app

RUN go build -o app ./cmd/main.go

FROM alpine

COPY --from=build /app /

CMD [ "/app" ]
FROM golang:1.15 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o contacts-api

FROM debian:buster-slim
COPY --from=builder /app/contacts-api /contacts-api
CMD [ "/contacts-api" ]
FROM golang:alpine AS build
WORKDIR /usr/src/app
COPY go.* ./
RUN go mod tidy
COPY . .
RUN GOOS=linux go build -o worker ./cmd/worker/main.go

FROM alpine:latest
COPY --from=build /usr/src/app/worker /usr/local/bin
CMD [ "/usr/local/bin/worker" ]
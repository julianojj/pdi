FROM golang:alpine AS build
WORKDIR /usr/src/app
COPY go.* ./
RUN go mod tidy
COPY . .
RUN GOOS=linux go build -o api ./main.go

FROM alpine:latest
COPY --from=build /usr/src/app/api /usr/local/bin
CMD [ "/usr/local/bin/api" ]

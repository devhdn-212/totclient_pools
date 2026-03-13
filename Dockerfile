FROM golang:alpine AS goqu
WORKDIR /go/src/github.com/devhdn-212/master
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as goqurelease
WORKDIR /app
RUN apk add tzdata
RUN mkdir -p ./frontend/public
COPY --from=goqu /go/src/github.com/devhdn-212/github.com/master/app .
COPY --from=goqu /go/src/github.com/devhdn-212/github.com/master/.env /app/.env

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 1010
CMD ["./app"]
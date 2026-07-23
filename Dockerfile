FROM golang:alpine AS totmodern_apiclient
WORKDIR /go/src/github.com/devhdn-212/totclient_api
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as totmodern_apiclient_release
WORKDIR /app
RUN apk add tzdata
COPY --from=totmodern_apiclient /go/src/github.com/devhdn-212/totclient_api/app .
COPY --from=totmodern_apiclient /go/src/github.com/devhdn-212/totclient_api/.env /app/.env

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 6058
CMD ["./app"]
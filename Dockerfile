FROM golang:alpine AS totmodern_agen
WORKDIR /go/src/github.com/devhdn-212/totagen_api
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as totmodern_agen_release
WORKDIR /app
RUN apk add tzdata
COPY --from=totmodern_agen /go/src/github.com/devhdn-212/totagen_api/app .
COPY --from=totmodern_agen /go/src/github.com/devhdn-212/totagen_api/.env /app/.env

ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 6059
CMD ["./app"]
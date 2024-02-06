#FROM alpine:3.19.1
FROM busybox
WORKDIR /app

COPY pltapi .
COPY config.json .

EXPOSE 8080

CMD ["/app/pltapi"]
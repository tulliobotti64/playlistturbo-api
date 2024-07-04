FROM busybox
WORKDIR /app

COPY pltapi .
COPY config.json .

EXPOSE 6061

CMD ["/app/pltapi"]
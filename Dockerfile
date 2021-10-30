FROM alpine
COPY ./src/server/http/static /swaggerui
COPY ./data /data
COPY ./bin/api /
CMD ["/api"]

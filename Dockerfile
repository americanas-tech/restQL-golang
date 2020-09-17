FROM alpine:3.11

RUN mkdir -p /restql

COPY ./bin/restQL /restql/api

WORKDIR /restql

RUN chmod +x ./api

CMD "./api"

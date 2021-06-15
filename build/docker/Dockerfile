FROM ubuntu:20.04

WORKDIR /app

COPY ./bin .

RUN mkdir /app/files
RUN apt-get update
RUN apt-get install -y ca-certificates

EXPOSE 8080
CMD ["/app/server"]

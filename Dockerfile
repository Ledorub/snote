FROM golang:1.23.1 AS build-stage

WORKDIR /usr/local/src/snote

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/snote ./cmd/app

FROM scratch

WORKDIR /usr/local/bin

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /usr/local/bin/snote ./snote

EXPOSE 8080

ENTRYPOINT ["./snote"]
CMD ["--port 8080"]

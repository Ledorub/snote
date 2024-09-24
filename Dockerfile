FROM golang:1.23.1 AS build-stage

WORKDIR /usr/local/src/snote

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/snote ./cmd/app

FROM build-stage AS release-stage

RUN useradd app
WORKDIR /usr/local/bin

COPY --from=build-stage /usr/local/bin/snote ./snote

USER app:app
EXPOSE 8080

ENTRYPOINT ["./snote"]
CMD ["--port 8080"]
